package deje

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
	Events     ObjectManager
	Quorums    serial.QuorumSet
	Timestamps ObjectManager
}

// Create a new, blank Document, with fields initialized.
func NewDocument() Document {
	return Document{
		Events:     NewObjectManager(),
		Quorums:    make(serial.QuorumSet),
		Timestamps: NewObjectManager(),
	}
}

// Copies the data from a DocumentFile into a Document.
func (d *Document) FromFile(df *serial.DocumentFile) {
	d.Channel = df.Channel
	d.Events = NewObjectManager()
	d.Quorums = df.Quorums

	for _, value := range df.Events {
		d.Events.Register(EventFromSerial(value))
	}
}

// Copies the data from a Document into a DocumentFile.
func (d *Document) ToFile() *serial.DocumentFile {
	df := &serial.DocumentFile{
		Channel: d.Channel,
		Events:  make(serial.EventSet),
		Quorums: d.Quorums,
	}

	for key, value := range d.Events.GetItems() {
		ev, ok := value.(Event)
		if ok {
			df.Events[key] = ev.ToSerial()
		}
	}

	return df
}
