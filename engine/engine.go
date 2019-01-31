package engine

import (
	"context"
	"time"
)

// Engine manages a list of queues and a continuous stream of events.
type Engine struct {
	state *DataStore
}

// New returns a new instance of the Bifrost Engine.
func New() *Engine {
	return &Engine{
		state: NewDataStore(),
	}
}

// Process will read and process events from eventCh until the context is
// done.
func (eng *Engine) Process(ctx context.Context, eventCh chan Event) error {
	go startSnapshots(ctx, eventCh)

	for {
		select {
		case evt := <-eventCh:
			eng.processEvent(evt)
		case <-ctx.Done():
			return nil
		}
	}
}

func (eng *Engine) processEvent(evt Event) {
	evt.Transition(eng.state)
}

func startSnapshots(ctx context.Context, eventCh chan Event) {
	timer := SnapshotTimer(
		ctx,
		10*time.Second,
		DefaultEncoderFactory,
		DefaultWriteCloserFactory,
	)

	for snapshotEvt := range timer {
		eventCh <- snapshotEvt
	}
}
