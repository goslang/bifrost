package engine

type Listener func(state *DataStore) Event

func ChangeListener(queueName string) (Listener, <-chan []byte) {
	// We're using lock to make sure only one publish is happening at a time.
	lock := make(chan bool, 1)
	unlock := func() { <-lock }
	publishCh := make(chan []byte)

	listener := func(state *DataStore) Event {
		select {
		case lock <- true:
		default:
			// Someone is reading this queue, nothing for us to do.
			return nil
		}

		q, ok := state.Buffers[queueName]
		if !ok || q.IsEmpty() {
			return nil
		}

		return Pop(queueName, func(message []byte) {
			defer unlock()
			publishCh <- message
		})
	}

	return listener, publishCh
}
