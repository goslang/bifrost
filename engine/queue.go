package engine

import (
	"bytes"
	"encoding/gob"
	"reflect"
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
		Size:   size,
		Buffer: make([][]byte, 0, size),
	}

	q.init()
	return q
}

// Copy returns a deep copy of the queue and it's messages.
func (q *Queue) Copy() *Queue {
	newQ := *q
	newQ.Buffer = make([][]byte, len(q.Buffer), q.Size)

	copy(newQ.Buffer, q.Buffer)

	return &newQ
}

func (q *Queue) init() {
	q.limiter = make(chan bool, q.Size)

	for range q.Buffer {
		q.limiter <- true
	}
}

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

func (q *Queue) GobDecode(raw []byte) error {
	reader := bytes.NewReader(raw)
	decoder := gob.NewDecoder(reader)

	for _, val := range q.EncodableValues() {
		if err := decoder.DecodeValue(val); err != nil {
			return err
		}
	}

	q.init()
	return nil
}

func (q *Queue) GobEncode() ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := gob.NewEncoder(buffer)

	for _, val := range q.EncodableValues() {
		if err := encoder.EncodeValue(val); err != nil {
			return nil, err
		}
	}

	return buffer.Bytes(), nil
}

func (q *Queue) EncodableValues() []reflect.Value {
	return []reflect.Value{
		reflect.ValueOf(&q.Size),
		reflect.ValueOf(&q.Buffer),
	}
}
