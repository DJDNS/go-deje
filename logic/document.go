package logic

import (
	"github.com/campadrenalin/go-deje/manager"
	"github.com/campadrenalin/go-deje/model"
	"github.com/campadrenalin/go-deje/state"
)

// A document is a single managed DEJE object, associated with
// a single immutable topic string, and self-describing its
// actions and permissions as part of the content.
//
// The content of a Document is the result of applying the
// "official" chain of history, in much the same way that the
// Bitcoin ledger is the result of playing through the transactions
// in every block of the longest valid blockchain.
type Document struct {
	Topic      string
	State      *state.DocumentState
	Events     *manager.EventManager
	Quorums    *manager.QuorumManager
	Timestamps *manager.TimestampManager
}

// Create a new, blank Document, with fields initialized.
func NewDocument() Document {
	return Document{
		State:      state.NewDocumentState(),
		Events:     manager.NewEventManager(),
		Quorums:    manager.NewQuorumManager(),
		Timestamps: manager.NewTimestampManager(),
	}
}

// Copies the data from a DocumentFile into a Document.
func (d *Document) FromFile(df *model.DocumentFile) {
	d.Topic = df.Topic
	d.Events = manager.NewEventManager()
	d.Quorums = manager.NewQuorumManager()

	d.Events.DeserializeFrom(df.Events)
	d.Quorums.DeserializeFrom(df.Quorums)
}

// Copies the data from a Document into a DocumentFile.
func (d *Document) ToFile() *model.DocumentFile {
	df := &model.DocumentFile{
		Topic:   d.Topic,
		Events:  make(model.EventSet),
		Quorums: make(model.QuorumSet),
	}

	d.Events.SerializeTo(df.Events)
	d.Quorums.SerializeTo(df.Quorums)

	return df
}
