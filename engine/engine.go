package engine

import (
	"context"
	"fmt"
	"os"
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
	eng.restoreState()
	go eng.startSnapshots()

	for {
		select {
		case evt := <-eng.eventCh:
			evt.Transition(eng.state)
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

func (eng *Engine) restoreState() {
	reader, err := DefaultReadCloser(eng.conf.snapshotFilename)
	if err != nil {
		// TODO: This should introspect the error to determine whether it is
		// fatal or not.
		fmt.Fprintln(os.Stderr, "ERROR: Failed to load snapshot file!")
		return
	}

	// TODO: ignoring error
	decoder, _ := DefaultDecoder(reader)

	err = decoder.Decode(&eng.state)
	if err != nil {
		// TODO: This should be more graceful than a panic, but if we can't
		// decode the file it should be considered a fatal error.
		fmt.Fprintln(os.Stderr, "ERROR: Failed to decode snapshot file!")
		panic(err)
	}
}

func (eng *Engine) startSnapshots() {
	timer := SnapshotTimer(
		eng.conf.ctx,
		eng.conf.snapshotInterval,
		DefaultEncoder,
		DefaultWriteCloserFactory(eng.conf.snapshotFilename),
	)

	for snapshotEvt := range timer {
		eng.eventCh <- snapshotEvt
	}
}
