package engine

import (
	"context"
)

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

		if ok := q.push(message); ok {
			confirmCh <- ok
		}
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
		ds.Buffers[name] = NewQueue(5)
	}

	return fn
}

func RemoveChannel(name string) Event {
	var fn EventFn = func(ds *DataStore) {
		delete(ds.Buffers, name)
	}

	return fn
}

func Pop(ctx context.Context, queueName string) (Event, <-chan []byte) {
	publishCh := make(chan []byte)

	var fn EventFn = func(ds *DataStore) {
		q, ok := ds.Buffers[queueName]
		if !ok {
			close(publishCh)
			return
		}

		go func() {
			select {
			case message := <-q.pop():
				publishCh <- message
			case <-ctx.Done():
			}
			close(publishCh)
		}()
	}

	return fn, publishCh
}

func PopNow(queueName string) (Event, <-chan []byte) {
	publishCh := make(chan []byte)

	var fn EventFn = func(ds *DataStore) {
		defer close(publishCh)

		q, ok := ds.Buffers[queueName]
		if !ok {
			return
		}

		select {
		case message, ok := <-q.pop():
			if ok {
				publishCh <- message
			}
		default:
		}
	}

	return fn, publishCh
}
