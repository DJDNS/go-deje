package model

import "github.com/campadrenalin/go-deje/serial"

// A document is a single managed DEJE object, associated with
// a single immutable IRCLocation, and self-describing its
// actions and permissions as part of the content.
//
// The content of a Document is the result of applying the
// "official" chain of history, in much the same way that the
// Bitcoin ledger is the result of playing through the transactions
// in every block of the longest valid blockchain.
type Document struct {
	Channel    serial.IRCLocation
	Events     EventManager
	Quorums    QuorumManager
	Timestamps TimestampManager
}

// Create a new, blank Document, with fields initialized.
func NewDocument() Document {
	return Document{
		Events:     NewEventManager(),
		Quorums:    NewQuorumManager(),
		Timestamps: NewTimestampManager(),
	}
}

// Copies the data from a DocumentFile into a Document.
func (d *Document) FromFile(df *serial.DocumentFile) {
	d.Channel = df.Channel
	d.Events = NewEventManager()
	d.Quorums = NewQuorumManager()

	d.Events.DeserializeFrom(df.Events)
	d.Quorums.DeserializeFrom(df.Quorums)
}

// Copies the data from a Document into a DocumentFile.
func (d *Document) ToFile() *serial.DocumentFile {
	df := &serial.DocumentFile{
		Channel: d.Channel,
		Events:  make(serial.EventSet),
		Quorums: make(serial.QuorumSet),
	}

	d.Events.SerializeTo(df.Events)
	d.Quorums.SerializeTo(df.Quorums)

	return df
}
