package document

import (
	"encoding/json"
	"io"

	"github.com/DJDNS/go-deje/state"
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
	Topic string               `json:"topic"`
	State *state.DocumentState `json:"-"`

	// Do not modify the contents of the following fields!
	// They're there for you to have convenient and uninhibited
	// READ-ONLY access. If you try to add or remove things manually,
	// you run the risk of doing so inconsistently.
	//
	// Please just use the Thing.Register() and Thing.Unregister()
	// methods, and when it comes to these fields, LOOK BUT DON'T TOUCH.
	Events         EventSet             `json:"events"`
	EventsByParent map[string]EventSet  `json:"-"`
	Quorums        QuorumSet            `json:"quorums"`
	QuorumsByEvent map[string]QuorumSet `json:"-"`
	Timestamps     []string             `json:"timestamps"`
}

// Create a new, blank Document, with fields initialized.
func NewDocument() Document {
	return Document{
		State:          state.NewDocumentState(),
		Events:         make(EventSet),
		EventsByParent: make(map[string]EventSet),
		Quorums:        make(QuorumSet),
		QuorumsByEvent: make(map[string]QuorumSet),
		Timestamps:     make([]string, 0),
	}
}

// Serialize JSON data to an io.Writer.
func (doc *Document) Serialize(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(doc)
}

// Deserialize JSON data from an io.Reader.
func (doc *Document) Deserialize(r io.Reader) error {
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(doc); err != nil {
		return err
	}

	// Copy Events to avoid clobbering when we fix keys
	var index int
	events_copy := make([]Event, len(doc.Events))
	for _, item := range doc.Events {
		events_copy[index] = *item
		index++
	}
	doc.Events = make(EventSet)

	// Same for Quorums
	index = 0
	quorums_copy := make([]Quorum, len(doc.Quorums))
	for _, item := range doc.Quorums {
		quorums_copy[index] = *item
		index++
	}
	doc.Quorums = make(QuorumSet)

	// Integrate through registration
	for i := range events_copy {
		events_copy[i].Doc = doc
		events_copy[i].Register()
	}
	for i := range quorums_copy {
		quorums_copy[i].Doc = doc
		quorums_copy[i].Register()
	}
	return nil
}
