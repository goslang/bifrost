package engine

// An event is anything that changes the application's state.
type Event interface {
	Transition(*DataStore)
}

type EventFn func(*DataStore)

func (fn EventFn) Transition(ds *DataStore) {
	fn(ds)
}

func PushMessage(name string, message []byte) (Event, <-chan bool) {
	confirmCh := make(chan bool, 1)

	var fn EventFn = func(ds *DataStore) {
		defer close(confirmCh)

		q, ok := ds.Buffers[name]
		if !ok {
			return
		}

		newQ := q.push(message)
		ds.Buffers[name] = newQ
		confirmCh <- true
	}

	return fn, confirmCh
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
