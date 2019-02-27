package engine

import (
	"context"
	"os"
	"sync"
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

	listenerIDMu      sync.Mutex
	currentListenerID int64
	listeners         map[int64]ListenerPair
}

// New returns a new instance of the Bifrost Engine.
func New(opts ...Opt) *Engine {
	eng := &Engine{
		state:   NewDataStore(),
		eventCh: make(chan Event),
	}

	//eng.Process(node.Events())

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
func (eng *Engine) Run() error {
	if err := eng.restoreState(); err != nil {
		return err
	}

	eng.startSnapshots()

	for {
		select {
		case evt := <-eng.eventCh:
			eng.applyEvent(evt)
		case <-eng.conf.ctx.Done():
			return nil
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

func (eng *Engine) applyEvent(evt Event) {
	changeSet := evt.Transition(eng.state)

	// TODO: This is not an efficient way to kick off the matching listeners.
	for _, pair := range eng.listeners {
		if pair.Matcher.Match(evt) {
			// Listeners are safe to run on their own goroutines.
			go pair.Listener(evt, changeSet)
		}
	}
}

func (eng *Engine) ListQueues() []QueueDetail {
	var details []QueueDetail

	for name, q := range eng.state.Buffers {
		detail := QueueDetail{
			Name: name,
			Max:  uint(q.Size),
			Size: uint(len(q.Buffer)),
		}

		details = append(details, detail)
	}

	return details
}

func (eng *Engine) GetQueueDetails(name string) (QueueDetail, bool) {
	q, ok := eng.state.Buffers[name]
	if !ok {
		return QueueDetail{}, false
	}

	detail := QueueDetail{
		Name: name,
		Max:  uint(q.Size),
		Size: uint(len(q.Buffer)),
	}

	return detail, true
}

// Stats returns an interface for querying the engine about it's current
// statistics.
func (eng *Engine) StatsAPI() StatsAPI {
	// Engine implements the StatsAPI, so just return it to expose access to
	// the stats subset of the API.
	return eng
}

// Listener returns an interface that encapsulates the listener portion of the
// Engine API.
func (eng *Engine) ListenerAPI() ListenerAPI {
	return eng
}

// Register registers an event listener with the engine and returns its
// listenerID.
func (eng *Engine) Register(matcher EventMatcher, l Listener) int64 {
	eng.listenerIDMu.Lock()
	defer eng.listenerIDMu.Unlock()

	eng.currentListenerID++
	eng.listeners[eng.currentListenerID] = ListenerPair{
		Listener: l,
		Matcher:  matcher,
	}

	return eng.currentListenerID
}

// Deregister removes the Listener with the corresponding listenerID.
func (eng *Engine) Deregister(listenerID int64) {
	delete(eng.listeners, listenerID)
}

func (eng *Engine) restoreState() error {
	reader, err := DefaultReadCloser(eng.conf.snapshotFilename)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	decoder, err := DefaultDecoder(reader)
	if err != nil {
		return err
	}

	err = decoder.Decode(&eng.state)
	if err != nil {
		return err
	}

	return nil
}

func (eng *Engine) startSnapshots() {
	timer := SnapshotTimer(
		eng.conf.ctx,
		eng.conf.snapshotInterval,
		DefaultEncoder,
		DefaultWriteCloserFactory(eng.conf.snapshotFilename),
	)

	eng.Process(eng.conf.ctx, timer)
}
