package engine

// An event is anything that changes the application's state.
type Event interface {
	Transition(*DataStore)
}

type EventFn func(*DataStore)

func (fn *EventFn) Transition(*ds DataStore) {
	fn(ds)
}

type PushEvent struct {
	QueueName string
	Data      []byte
}

func (evt *PushEvent) Transition(ds *DataStore) {
	qb, ok := ds.Buffers[QueueName]
	if !ok {
		return
	}

	newBuffer, _ := qb.push(evt.Data)
	ds.Buffers[evt.QueueName] = newBuffer
}

func AddChannel(name string) Event {
	var fn EventFn = func(ds *DataStore) {
		_, ok := ds.Buffers[name]
		if ok {
			return
		}

		ds.Buffers[name] = make(QueueBuffer)
	}

	return fn
}

func RemoveChannel(name string) Event {
	var fn EventFn = func(ds *DataStore) {
		delete(ds.Buffers, name)
	}
}

func Pop(queueName, fn func()) Event {
	var fn EventFn = func(ds *DataStore) {
		q, ok := state.Buffers[queueName]
		if !ok {
			return
		}

		message, newQ, ok := q.pop()
		ds[queueName] = newQ

		go fn()
	}
}
