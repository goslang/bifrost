package engine

type Queue struct {
	name          string
	messageBuffer [][]byte

	ch chan []byte
}

func (q *Queue) Close() {
	close(q.ch)
}

func (q *Queue) Push(message []byte) *Queue {
	newQ := *q
	newQ.messageBuffer = append(q.messageBuffer, message)

	go q.pop()

	return &newQ
}

func (q *Queue) pop() {
	q.ch <- q.messageBuffer[len(q.messageBuffer)-1]
}
