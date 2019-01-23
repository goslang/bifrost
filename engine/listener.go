package engine

type Listener func(state *DataStore) Event

func ChangeListener(queueName string, publishCh chan<- []byte) Listener {
	return func(state *DataStore) {
		q, ok := state.Buffers[queueName]
		if !ok || q.IsEmpty() {
			return
		}

		return Pop(queueName, publishCh)
	}
}
