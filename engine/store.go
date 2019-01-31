package engine

type DataStore struct {
	Buffers map[string]*Queue
}

func NewDataStore() *DataStore {
	return &DataStore{
		Buffers: make(map[string]*Queue),
	}
}

func (ds *DataStore) Copy() *DataStore {
	newDs := NewDataStore()

	for name, buf := range ds.Buffers {
		newDs.Buffers[name] = buf.Copy()
	}

	return newDs
}
