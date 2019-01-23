package engine

type DataStore struct {
	Buffers map[string]QueueBuffer

	// TODO: Currently this will form an endless chain of previous states,
	// which, while really cool, will eventually consume all memory and
	// thoughts of loved ones gone before... so we need to put a cap on this.
	prev *DataStore
	next *DataStore
}

func NewDataStore() *DataStore {
	return &DataStore{
		Buffers: make(map[string]QueueBuffer,
	}
}

func (ds *DataStore) MakeNewState() *DataStore {
	newDs := NewDataStore()

	newDs.prev = ds
	ds.next = newDs

	for k, v := range ds.Buffers {
		newDs.Buffers[k] = v
	}

	return newDs
}
