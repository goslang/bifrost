package engine

import (
	"sync"
)

type Queue struct {
	Buffer [][]byte

	idx int
	mu  *sync.Mutex
	ch  chan []byte
}

func NewQueue() *Queue {
	q := &Queue{
		// TODO: Make this configurable.
		Buffer: make([][]byte, 5),
		mu:     &sync.Mutex{},
		ch:     make(chan []byte),
	}

	return q
}

func (q *Queue) Close() {
	close(q.ch)
}

func (q *Queue) push(message []byte) *Queue {
	newQ := *q
	newQ.Buffer = make([][]byte, len(q.Buffer))
	copy(newQ.Buffer, q.Buffer)

	println(">>>> Adding message to buffer")
	println(">>>>", string(message))
	newQ.write(message)
	println("idx =", q.idx)

	// TODO: Cap the maximum number of goroutines here.
	go func(message []byte) {
		newQ.ch <- message
	}(newQ.Buffer[newQ.readIdx()])

	return &newQ
}

func (q *Queue) pop() <-chan []byte {
	return q.ch
}

func (q *Queue) write(message []byte) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.Buffer[q.idx] = message
	q.idx++

	if q.idx == len(q.Buffer) {
		q.idx = 0
	}
}

func (q *Queue) readIdx() int {
	if q.idx == 0 {
		return len(q.Buffer) - 1
	}

	return q.idx - 1
}
