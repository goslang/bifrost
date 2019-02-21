package engine

// StatsAPI provides methods to access up-to-date statistics about the
// Engine's current state.
type StatsAPI interface {
	GetQueueDetails(string) (QueueDetail, bool)
	ListQueues() []QueueDetail
}

// QueueDetail represents current stats for a particular queue.
type QueueDetail struct {
	Name string `json:"name"`
	Size uint   `json:"size"`
	Max  uint   `json:"max"`
}
