package engine

// Event is an action that causes the engine to transition from one state to
// the next.
type Event interface {
	Transition(*DataStore) ChangeSet
}

// EventFn implement the Event interface, allowing you to use simple functions
// instead of objects as events.
type EventFn func(*DataStore) ChangeSet

// Transition Invokes the function.
func (fn EventFn) Transition(ds *DataStore) ChangeSet {
	return fn(ds)
}

// ChangeSet represents the change of state caused by a single transition.
// Listeners will receive a copy of the ChangeSet after the transition has
// been applied.
type ChangeSet struct {
	Pushed []byte
	Popped []byte

	Added   string
	Removed string
}

type Push struct {
	Channel string
	Message []byte
}

func (evt Push) Transition(store *DataStore) ChangeSet {
	q, ok := store.Buffers[evt.Channel]
	if !ok {
		panic("This is not OK")
	}

	if ok := q.push(evt.Message); !ok {
		panic("It's still not OK")
	}

	return ChangeSet{Pushed: evt.Message}
}

type Pop struct {
	Channel  string
	ClientID string
}

func (evt Pop) Transition(store *DataStore) ChangeSet {
	q, ok := store.Buffers[evt.Channel]
	if !ok {
		return ChangeSet{}
	}

	message, ok := q.pop()
	if !ok {
		return ChangeSet{}
	}

	return ChangeSet{Popped: message}
}
