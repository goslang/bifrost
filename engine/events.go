package engine

import (
	"context"
)

// An event is anything that changes the application's state.
type Event interface {
	Transition(*DataStore)
}

// EventFn implement the Event interface, allowing you to use simple functions
// instead of objects as events.
type EventFn func(*DataStore)

// Transition Invokes the function.
func (fn EventFn) Transition(ds *DataStore) {
	fn(ds)
}

// PushMessage creates an event that sends `message` to the channel named
// `name`. It also returns a confirmation channel that it will write to once
// the message has been successfully added to the queue. If the message cannot
// be added, the confirmation channel will be closed without sending a value.
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

// AddChannel returns an Event that will create a new channel named `name`.
func AddChannel(name string, size int) Event {
	var fn EventFn = func(ds *DataStore) {
		_, ok := ds.Buffers[name]
		if ok {
			return
		}

		ds.Buffers[name] = NewQueue(size)
	}

	return fn
}

// RemoveChannel returns an Event that removes the channel named `name` from
// the engine.
func RemoveChannel(name string) Event {
	var fn EventFn = func(ds *DataStore) {
		delete(ds.Buffers, name)
	}

	return fn
}

// Pop returns an Event that will read the next message from the specified
// queue and write it to the returned channel.
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

// PopNow returns an event that will pop the next item off the queue, but fails
// immediately if no messages are preseent to be read.
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
