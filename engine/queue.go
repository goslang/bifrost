package engine

import (
	"sync"
)

// Queue manages the state for a queue of messages.
type Queue struct {
	// Buffer is an array of arbitrary messages.
	Buffer [][]byte

	// Size is the maximum number of messages  to buffer in this queue at any
	// given time.
	Size uint

	// limiterCh guarantees that we will never A) push onto a full queue, or
	// B) pop off of an empty queue. Don't use this directly, instead call
	// `q.limiter()` to ensure that it has been properly initialized.
	limiterCh chan bool

	// mu protects concurrent operations on the Queue.
	mu sync.Mutex

	// init ensures that the limiterCh is initialized once and only once.
	init sync.Once
}

// NewQueue creates a Queue that contains up to `size` messages.
func NewQueue(size uint) *Queue {
	q := &Queue{
		Size:   size,
		Buffer: make([][]byte, 0, size),
	}

	return q
}

// Copy returns a deep copy of the queue and it's messages.
func (q *Queue) Copy() *Queue {
	newQ := *q
	newQ.Buffer = make([][]byte, len(q.Buffer), q.Size)

	copy(newQ.Buffer, q.Buffer)

	return &newQ
}

func (q *Queue) push(message []byte) bool {
	select {
	case q.limiter() <- true:
	default:
		return false
	}

	q.safePush(message)
	return true
}

// pop pops the next item off of the queue. Its second return value is set to
// false if no data is currently available.
func (q *Queue) pop() ([]byte, bool) {
	select {
	case <-q.limiter():
	default:
		return nil, false
	}

	return q.safePop(), true
}

// listenOne pops the next item off of the queue and sends it to the returned
// channel when data becomes available.
func (q *Queue) listenOne() <-chan []byte {
	ch := make(chan []byte)

	go func() {
		defer close(ch)

		message := q.safePop()
		ch <- message
	}()

	return ch
}

func (q *Queue) safePop() []byte {
	q.mu.Lock()
	defer q.mu.Unlock()

	message := q.Buffer[0]
	q.Buffer = q.Buffer[1:]
	return message
}

func (q *Queue) safePush(message []byte) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.Buffer = append(q.Buffer, message)
}

// limiter wraps the q's limiterCh to make sure it is initialized properly.
// This way we can guarantee that listeners will be setup properly even if the
// queue has just been reloaded from disk.
func (q *Queue) limiter() chan bool {
	q.init.Do(func() {
		q.limiterCh = make(chan bool, q.Size)

		for range q.Buffer {
			q.limiterCh <- true
		}
	})

	return q.limiterCh
}
