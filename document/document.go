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
