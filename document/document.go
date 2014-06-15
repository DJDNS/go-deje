package document

import "github.com/campadrenalin/go-deje/state"

// A document is a single managed DEJE object, associated with
// a single immutable topic string, and self-describing its
// actions and permissions as part of the content.
//
// The content of a Document is the result of applying the
// "official" chain of history, in much the same way that the
// Bitcoin ledger is the result of playing through the transactions
// in every block of the longest valid blockchain.
type Document struct {
	Topic string
	State *state.DocumentState

	// Do not modify the contents of the following fields!
	// They're there for you to have convenient and uninhibited
	// READ-ONLY access. If you try to add or remove things manually,
	// you run the risk of doing so inconsistently.
	//
	// Please just use the Thing.Register() and Thing.Unregister()
	// methods, and when it comes to these fields, LOOK BUT DON'T TOUCH.
	Events         EventSet
	EventsByParent map[string]EventSet
	Quorums        QuorumSet
	QuorumsByEvent map[string]QuorumSet
}

// Create a new, blank Document, with fields initialized.
func NewDocument() Document {
	return Document{
		State:          state.NewDocumentState(),
		Events:         make(EventSet),
		EventsByParent: make(map[string]EventSet),
		Quorums:        make(QuorumSet),
		QuorumsByEvent: make(map[string]QuorumSet),
	}
}

// Copies the data from a DocumentFile into a Document.
func (d *Document) FromFile(df *DocumentFile) {
	d.Topic = df.Topic
	d.Events = make(EventSet)
	d.EventsByParent = make(map[string]EventSet)
	d.Quorums = make(QuorumSet)
	d.QuorumsByEvent = make(map[string]QuorumSet)

	for _, ev := range df.Events {
		ev.Doc = d
		ev.Register()
	}
	for _, q := range df.Quorums {
		q.Doc = d
		q.Register()
	}
}

// Copies the data from a Document into a DocumentFile.
func (d *Document) ToFile() *DocumentFile {
	df := &DocumentFile{
		Topic:   d.Topic,
		Events:  make(EventSet),
		Quorums: make(QuorumSet),
	}

	for key, ev := range d.Events {
		df.Events[key] = ev
	}
	for key, q := range d.Quorums {
		df.Quorums[key] = q
	}

	return df
}
