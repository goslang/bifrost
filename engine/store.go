package engine

type DataStore struct {
	Buffers map[string]*Queue

	// TODO: Currently this will form an endless chain of previous states,
	// which, while really cool, will eventually consume all memory and
	// thoughts of loved ones gone before... so we need to put a cap on this.
	prev *DataStore
	next *DataStore
}

func NewDataStore() *DataStore {
	return &DataStore{
		Buffers: make(map[string]*Queue),
	}
}
