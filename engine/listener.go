package engine

type ListenerAPI interface {
	Register(EventMatcher, Listener) (listenerId int64)
	Deregister(listenerId int64)
}

type ListenerPair struct {
	Listener Listener
	Matcher  EventMatcher
}

type EventMatcher struct {
	ClientID    string
	EventType   string
	ChannelName string
}

func (em EventMatcher) Match(evt Event) (matches bool) {
	defer func() {
		if err := recover(); err != nil {
			matches = false
		}
	}()

	em.matchClient(evt)
	em.matchType(evt)
	em.matchChannel(evt)

	return true
}

func (em EventMatcher) matchClient(evt Event) {
	if em.ClientID == "" {
		return
	}

	pop, ok := evt.(Pop)
	if ok && pop.ClientID != em.ClientID {
		panic("ClientID does not match")
	}
}

func (em EventMatcher) matchType(evt Event) {
}

func (em EventMatcher) matchChannel(evt Event) {
}

type Listener func(Event, ChangeSet)

// PushListener will wait for matching push events and send a Pop event when
// they occur. To receive Popped items, setup a corresponding PopListener to
// receive them.
func PushListener(clientId string) (Listener, chan Event) {
	eventsCh := make(chan Event)

	fn := func(evt Event, _ ChangeSet) {
		push, ok := evt.(Push)
		if !ok {
			return
		}

		eventsCh <- Pop{
			Channel:  push.Channel,
			ClientID: clientId,
		}
	}

	return fn, eventsCh
}

// PopListener pops a single item off of the queue and sends it to the
// returned channel.
func PopListener() (Listener, <-chan []byte) {
	// I prefer classic rock myself

	publishCh := make(chan []byte)

	l := func(evt Event, changes ChangeSet) {
		message := changes.Popped
		publishCh <- message
		close(publishCh)
	}

	return l, publishCh
}
