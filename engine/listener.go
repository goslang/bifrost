package engine

type Listener func(state *DataStore)

func ChangeListener(queueName string, publishCh <-chan []byte) Listener {
	// We're using lock to make sure only one publish is happening at a time.
	//lock := make(chan bool, 1)
	//unlock := func() { <-lock }

	listener := func(state *DataStore) {
		//select {
		//case lock <- true:
		//default:
		//	return
		//}

		//q, ok := state.Buffers[queueName]
		//if !ok || q.IsEmpty() {
		//	return
		//}

		//go func() {
		//	defer unlock()
		//	message := <-popCh
		//	publishCh <- message
		//}()

		return
	}

	return listener
}
