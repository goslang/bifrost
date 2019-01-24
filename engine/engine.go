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
	state *DataStore
}

func New() *Engine {
	return &Engine{
		state: NewDataStore(),
	}
}

// Process will read and process events from eventCh until the context is
// done.
func (eng *Engine) Process(ctx context.Context, eventCh <-chan Event) error {
	println(">>>> Engine beginning to process")
	for {
		select {
		case evt := <-eventCh:
			println(">>>> Processing event")
			eng.processEvent(evt)
		case <-ctx.Done():
			return nil
		}
	}
}

func (eng *Engine) processEvent(evt Event) {
	eng.state = eng.state.MakeNewState()
	evt.Transition(eng.state)
}
