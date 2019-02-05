package engine

import (
	"sync"
)

// Queue manages the state for a queue of messages.
type Queue struct {
	Buffer  []Message
	limiter chan bool

	mu sync.Mutex
}

// Message contains the actual data that should be passed to the consumer, and
// tracks whether it has been delivered.
type Message struct {
	Data      []byte
	Delivered bool
}

// NewQueue creates a Queue that contains up to `size` messages.
func NewQueue(size int) *Queue {
	q := &Queue{
		Buffer:  make([]Message, 0, size),
		limiter: make(chan bool, size),
	}

	return q
}

// Close safely closes the Queue.
func (q *Queue) Close() {
	//close(q.ch)
}

// Copy returns a deep copy of the queue and it's messages.
func (q *Queue) Copy() *Queue {
	newQ := *q
	newQ.Buffer = make([]Message, 0, cap(q.Buffer))

	copy(newQ.Buffer, q.Buffer)

	return &newQ
}

func (q *Queue) push(data []byte) bool {
	select {
	case q.limiter <- true:
	default:
		return false
	}

	message := Message{Data: data}
	q.safePush(message)

	return true
}

// pop pops the next item off of the queue. Its second return value is set to
// false if no data is currently available.
func (q *Queue) pop() ([]byte, bool) {
	select {
	case <-q.limiter:
	default:
		return nil, false
	}

	message := q.safePop()
	return message.Data, true
}

// listenOne pops the next item off of the queue and sends it to the returned
// channel when data becomes available.
func (q *Queue) listenOne() <-chan []byte {
	ch := make(chan []byte)

	go func() {
		defer close(ch)

		<-q.limiter
		message := q.safePop()

		ch <- message.Data
	}()

	return ch
}

func (q *Queue) safePop() Message {
	defer q.lockAndUnlock("safePop")()

	message := q.Buffer[0]
	q.Buffer = q.Buffer[1:]
	return message
}

func (q *Queue) safePush(message Message) {
	defer q.lockAndUnlock("safePush")()

	q.Buffer = append(q.Buffer, message)
}

func (q *Queue) lockAndUnlock(hint string) func() {
	q.mu.Lock()
	println("locked", hint)

	return func() {
		q.mu.Unlock()
		println("Unlocked", hint)
	}
}
