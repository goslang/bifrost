package engine

import (
	"sync"
)

// Queue manages the state for a queue of messages.
type Queue struct {
	Buffer  [][]byte
	Size    uint
	limiter chan bool

	mu sync.Mutex
}

// NewQueue creates a Queue that contains up to `size` messages.
func NewQueue(size uint) *Queue {
	q := &Queue{
		Buffer:  make([][]byte, 0, size),
		limiter: make(chan bool, size),
	}

	return q
}

// Copy returns a deep copy of the queue and it's messages.
func (q *Queue) Copy() *Queue {
	newQ := *q
	newQ.Buffer = make([][]byte, 0, cap(q.Buffer))

	copy(newQ.Buffer, q.Buffer)

	return &newQ
}

func (q *Queue) init()

func (q *Queue) push(message []byte) bool {
	select {
	case q.limiter <- true:
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
	case <-q.limiter:
	default:
		return nil, false
	}

	message := q.safePop()
	return message, true
}

// listenOne pops the next item off of the queue and sends it to the returned
// channel when data becomes available.
func (q *Queue) listenOne() <-chan []byte {
	ch := make(chan []byte)

	go func() {
		defer close(ch)

		<-q.limiter
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

// QueueAlias is only used for decoding Queue objects.
type QueueAlias Queue

func (q *Queue) GobDecode(raw []byte) error {
	var qa QueueAlias

	reader := bytes.NewReader(raw)
	if err := gob.NewDecoder(reader).Decode(&qa); err != nil {
		return err
	}

	*q = qa
	q.init()
}
