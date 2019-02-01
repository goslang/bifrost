package engine

// Queue is essentially a ring buffer.
type Queue struct {
	Buffer   []Message
	WriteIdx int
	limiter  chan bool

	// push channel that new items are sent to.
	ch chan []byte
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
		Buffer:  make([]Message, size),
		limiter: make(chan bool, size),
		ch:      make(chan []byte),
	}

	return q
}

// Close safely closes the Queue.
func (q *Queue) Close() {
	close(q.ch)
}

// Copy returns a deep copy of the queue and it's messages.
func (q *Queue) Copy() *Queue {
	newQ := *q
	newQ.Buffer = make([]Message, len(q.Buffer))

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
	q.write(message)

	go func() {
		q.ch <- message.Data
		message.Delivered = true
		<-q.limiter
	}()

	return true
}

func (q *Queue) pop() <-chan []byte {
	return q.ch
}

func (q *Queue) write(message Message) {
	newIdx := q.incr(q.WriteIdx)

	q.Buffer[q.WriteIdx] = message
	q.WriteIdx = newIdx

	return
}

// incr increments idx by 1 unless it equals len(q.Buffer), and then restarts
// it at 0.
func (q *Queue) incr(idx int) int {
	if idx == len(q.Buffer)-1 {
		return 0
	}
	return idx + 1
}
