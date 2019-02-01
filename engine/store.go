package engine

// DataStore contains the Engine's current state.
type DataStore struct {
	Buffers map[string]*Queue
}

// NewDataStore initializes a new DataStore
func NewDataStore() *DataStore {
	return &DataStore{
		Buffers: make(map[string]*Queue),
	}
}

// Copy returns a new datastore that is a deep copy of the current datastore.
func (ds *DataStore) Copy() *DataStore {
	newDs := NewDataStore()

	for name, buf := range ds.Buffers {
		newDs.Buffers[name] = buf.Copy()
	}

	return newDs
}
