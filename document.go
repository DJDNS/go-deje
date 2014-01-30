package deje

type Document struct {
	Channel    IRCLocation
	Events     EventSet
	Syncs      SyncSet
	Timestamps TimestampSet `json:""`
}

func NewDocument() Document {
	return Document{
		Events: make(EventSet),
		Syncs:  make(SyncSet),
	}
}
