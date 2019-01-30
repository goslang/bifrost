package engine

import (
	"context"
	"errors"
	"time"
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
		DefaultWriterFactory,
	)

	for snapshotEvt := range timer {
		eventCh <- snapshotEvt
	}
}
