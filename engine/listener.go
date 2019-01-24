package engine

type Listener func(state *DataStore) Event

func ChangeListener(queueName string) (Listener, publishCh <-chan []byte) {
	// We're using lock to make sure only one publish is happening at a time.
	lock := make(chan bool, 1)
	unlock := func() { <-lock }
	publishCh := make(chan []byte)

	return func(state *DataStore) {
		select {
		case lock <- 1:
		default:
			// Someone is reading this queue, nothing for us to do.
			return
		}

		q, ok := state.Buffers[queueName]
		if !ok || q.IsEmpty() {
			return
		}

		return Pop(queueName, func(message []byte) {
			defer unlock()
			publishCh <- message
		})
	}
}
