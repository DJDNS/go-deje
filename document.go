package deje

type Document struct {
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
