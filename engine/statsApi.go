package engine

// StatsAPI provides methods to access up-to-date statistics about the
// Engine's current state.
type StatsAPI interface {
	GetChannelDetails(string) (ChannelDetail, bool)
	ListChannels() []ChannelDetail
}

// ChannelDetail represents current stats for a particular queue.
type ChannelDetail struct {
	Name string `json:"name"`
	Size uint   `json:"size"`
	Max  uint   `json:"max"`
}
