package deje

// A document is a single managed DEJE object, associated with
// a single immutable IRCLocation, and self-describing its
// actions and permissions as part of the content.
//
// The content of a Document is the result of applying the
// "official" chain of history, in much the same way that the
// Bitcoin ledger is the result of playing through the transactions
// in every block of the longest valid blockchain.
type Document struct {
	Channel    IRCLocation
	Events     EventSet
	Syncs      SyncSet
	Timestamps TimestampManager
}

// Used for serializing and deserializing docs to files.
//
// This allows us to use more complicated structures for actual
// documents, that allow for storing Timestamps, and other data
// that we must not trust the file to provide.
type DocumentFile struct {
	Channel IRCLocation
	Events  EventSet
	Syncs   SyncSet
}

// Create a new, blank Document, with fields initialized.
func NewDocument() Document {
	return Document{
		Events: make(EventSet),
		Syncs:  make(SyncSet),
	}
}

// Copies the data from a DocumentFile into a Document.
func (d *Document) FromFile(df *DocumentFile) {
	d.Channel = df.Channel
	d.Events = df.Events
	d.Syncs = df.Syncs
}

// Copies the data from a Document into a DocumentFile.
func (d *Document) ToFile() *DocumentFile {
	return &DocumentFile{
		Channel: d.Channel,
		Events:  d.Events,
		Syncs:   d.Syncs,
	}
}
