package engine

type QueueBuffer [][]byte

func (qb *QueueBuffer) pop() ([]byte, *QueueBuffer, ok) {
	if qb.IsEmpty() {
		return nil, *qb, false
	}

	message := qb[0]
	newQ := append(qb[:0:0], qb[1:len(qb)])

	return message, newQ, true
}

func (qb *QueueBuffer) push(message []byte) (*QueueBuffer, ok) {
	if qb.IsFull() {
		return qb, false
	}

	newQ := append(qb, message)
	return newQ
}

func (qb *QueueBuffer) IsFull() bool {
	return len(qb) == cap(qb)
}

func (qb *QueueByffer) IsEmpty() bool {
	return len(qb) < 1
}
