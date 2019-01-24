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
	q, ok := ds.Buffers[evt.QueueName]
	if !ok {
		return
	}

	newQ := q.push(evt.Data)
	ds.Buffers[evt.QueueName] = newQ
}

func AddChannel(name string) Event {
	var fn EventFn = func(ds *DataStore) {
		_, ok := ds.Buffers[name]
		if ok {
			return
		}

		// TODO: Configure the buffer size.
		ds.Buffers[name] = NewQueue()
	}

	return fn
}

func RemoveChannel(name string) Event {
	var fn EventFn = func(ds *DataStore) {
		delete(ds.Buffers, name)
	}

	return fn
}

func Pop(queueName string) (Event, <-chan []byte) {
	publishCh := make(chan []byte)

	var fn EventFn = func(ds *DataStore) {
		q, ok := ds.Buffers[queueName]
		if !ok {
			close(publishCh)
			return
		}

		go func() {
			message := <-q.pop()
			publishCh <- message
			close(publishCh)
		}()
	}

	return fn, publishCh
}
