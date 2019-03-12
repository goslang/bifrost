package engine

// DataStore contains the Engine's current state.
type DataStore struct {
	Channels map[string]*Channel
}

// NewDataStore initializes a new DataStore
func NewDataStore() *DataStore {
	return &DataStore{
		Channels: make(map[string]*Channel),
	}
}

// Copy returns a new datastore that is a deep copy of the current datastore.
func (ds *DataStore) Copy() *DataStore {
	newDs := NewDataStore()

	for name, buf := range ds.Channels {
		newDs.Channels[name] = buf.Copy()
	}

	return newDs
}
