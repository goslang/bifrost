package engine

import (
	"context"
	"errors"
)

var (
	ErrNoQueue     = errors.New("Queue not found.")
	ErrBufferFull  = errors.New("Queue buffer full.")
	ErrBufferEmpty = errors.New("Queue buffer empty.")
)

type Engine struct {
	state     *DataStore
	listeners map[int64]Listener

	// listenerID tracks the current to make sure it returns a new ID when a
	// new listener is added.
	listenerID int64
}

func New() *Engine {
	return &Engine{
		state:     NewDataStore(),
		listeners: make(map[int64]Listener),
	}
}

// Register Registers a new listener that will be called when the state
// changes. Returns an int64 ID that can be used to deregister the listener
// when it's no longer needed
func (eng *Engine) Register(listener Listener) int64 {
	nextId := eng.listenerID++
	eng.listeners[nextId] = listener
	return nextId
}

// Deregister will remove the listener with the matching ID, so it will not be
// run any more.
func (eng *Engine) Deregister(listenerID int64) {
	delete(eng.listeners, listenerID)
}

// Process will read and process events from eventCh until the context is
// done.
func (eng *Engine) Process(ctx context.Context, eventCh <-chan Event) error {
	for {
		select {
		case evt := <-evtCh:
			eng.processEvent(evt)
			eng.runListeners()
		case <-ctx.Done():
			return
		}
	}
}

func (eng *Engine) processEvt(evt Event) {
	eng.State = eng.State.MakeNewState()
	evt.Transition(eng.State)
}

func (eng *Engine) runListeners() {
	for _, listener := range eng.Listeners {
		go func(listener Listener) {
			listener(eng.State)
		}(listener)
	}
}

func (eng *Engine) Pop(channelName string) ([]byte, error) {
	queue, ok := eng.queues[channelName]
	if !ok {
		return nil, ErrNoQueue
	}

	select {
	case message := <-queue.ch:
		return message, nil
	default:
		return nil, ErrBufferEmpty
	}
}
