package engine

// Queue is essentially a ring buffer.
type Queue struct {
	Buffer   [][]byte
	writeIdx int
	limiter  chan bool

	// push channel that new items are sent to.
	ch chan []byte
}

func NewQueue(size int) *Queue {
	q := &Queue{
		Buffer:  make([][]byte, size),
		limiter: make(chan bool, size),
		ch:      make(chan []byte),
	}

	return q
}

func (q *Queue) Close() {
	close(q.ch)
}

func (q *Queue) push(message []byte) bool {
	select {
	case q.limiter <- true:
	default:
		return false
	}

	q.write(message)

	go func() {
		q.ch <- message
		<-q.limiter
	}()

	return true
}

func (q *Queue) pop() <-chan []byte {
	return q.ch
}

func (q *Queue) write(message []byte) {
	newIdx := q.incr(q.writeIdx)

	q.Buffer[q.writeIdx] = message
	q.writeIdx = newIdx

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
