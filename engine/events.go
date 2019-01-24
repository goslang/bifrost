package engine

// An event is anything that changes the application's state.
type Event interface {
	Transition(*DataStore)
}

type EventFn func(*DataStore)

func (fn EventFn) Transition(ds *DataStore) {
	fn(ds)
}

type PushEvent struct {
	QueueName string
	Data      []byte
}

func (evt *PushEvent) Transition(ds *DataStore) {
	qb, ok := ds.Buffers[evt.QueueName]
	if !ok {
		return
	}

	newQ, _ := qb.push(evt.Data)
	ds.Buffers[evt.QueueName] = newQ
}

func AddChannel(name string) Event {
	var fn EventFn = func(ds *DataStore) {
		_, ok := ds.Buffers[name]
		if ok {
			return
		}

		// TODO: Configure the buffer size.
		ds.Buffers[name] = make(QueueBuffer, 5)
	}

	return fn
}

func RemoveChannel(name string) Event {
	var fn EventFn = func(ds *DataStore) {
		delete(ds.Buffers, name)
	}

	return fn
}

func Pop(queueName string, callback func(message []byte)) Event {
	// TODO: accepting a `callback` is probably an antipattern

	var fn EventFn = func(ds *DataStore) {
		q, ok := ds.Buffers[queueName]
		if !ok {
			return
		}

		message, newQ, ok := q.pop()
		ds.Buffers[queueName] = newQ

		go callback(message)
	}

	return fn
}
