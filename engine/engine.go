package engine

import (
	"context"
	"time"
)

var defaultOpts = []Opt{
	Context(context.Background()),
	SnapshotInterval(10 * time.Second),
	SnapshotFilename("/usr/local/var/bifrost/snapshot.data"),
}

// Engine manages a list of queues and a continuous stream of events.
type Engine struct {
	state   *DataStore
	eventCh chan Event
	conf    Config
}

// New returns a new instance of the Bifrost Engine.
func New(opts ...Opt) *Engine {
	eng := &Engine{
		state:   NewDataStore(),
		eventCh: make(chan Event),
	}

	return eng.With(defaultOpts...).With(opts...)
}

// With applies the list of Opts to the engine. It returns the engine to make
// it chainable.
func (eng *Engine) With(opts ...Opt) *Engine {
	for _, opt := range opts {
		eng.conf = opt(eng.conf)
	}

	return eng
}

// Run starts the engine processing messages
func (eng *Engine) Run() {
	go eng.startSnapshots(eng.conf.ctx, eng.eventCh)

	for {
		select {
		case evt := <-eng.eventCh:
			eng.processEvent(evt)
		case <-eng.conf.ctx.Done():
			return
		}
	}
}

// Process multiplexes the provided eventCh onto the Engine's main event
// channel until ctx is finished.
func (eng *Engine) Process(ctx context.Context, eventCh <-chan Event) {
	go func() {
		for {
			select {
			case evt := <-eventCh:
				eng.eventCh <- evt
			case <-ctx.Done():
				return
			case <-eng.conf.ctx.Done():
				return
			}
		}
	}()
}

func (eng *Engine) processEvent(evt Event) {
	evt.Transition(eng.state)
}

func (eng *Engine) startSnapshots(ctx context.Context, eventCh chan Event) {
	timer := SnapshotTimer(
		ctx,
		eng.conf.snapshotInterval,
		DefaultEncoderFactory,
		DefaultWriteCloserFactory(eng.conf.snapshotFilename),
	)

	for snapshotEvt := range timer {
		eventCh <- snapshotEvt
	}
}
