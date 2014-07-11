package document

import "github.com/DJDNS/go-deje/util"

// Represents a complete set of approvals for an event.
// Quorums act as bridges between events and timestamps,
// indicating that an event was both common knowledge and
// considered a valid event chain (among others) at one
// time (the timestamp provides the time information).
type Quorum struct {
	Doc        *Document         `json:"-"`
	EventHash  string            `json:"event_hash"`
	Signatures map[string]string `json:"sigs"`
}

type QuorumSet map[string]*Quorum

func NewQuorum(evhash string) Quorum {
	return Quorum{
		EventHash:  evhash,
		Signatures: make(map[string]string),
	}
}

func (doc *Document) NewQuorum(event_hash string) Quorum {
	q := NewQuorum(event_hash)
	q.Doc = doc
	return q
}

func (q Quorum) GetKey() string {
	return q.Hash()
}
func (q Quorum) GetGroupKey() string {
	return q.EventHash
}
func (q Quorum) Eq(other Quorum) bool {
	return q.Hash() == other.Hash()
}

// Get the hash of the Quorum object.
func (q Quorum) Hash() string {
	hash, _ := util.HashObject(q)
	return hash
}

// Register with the Doc. This stores it in a hash-based location,
// so do not make changes to an Event after it has been registered.
func (q *Quorum) Register() {
	key := q.GetKey()
	q.Doc.Quorums[key] = q

	group_key := q.GetGroupKey()
	group, ok := q.Doc.QuorumsByEvent[group_key]
	if !ok {
		group = make(QuorumSet)
		q.Doc.QuorumsByEvent[group_key] = group
	}
	group[key] = q
}

// Unregister from the Doc. This also cleans up empty groups.
func (q *Quorum) Unregister() {
	key := q.GetKey()
	delete(q.Doc.Quorums, key)

	group_key := q.GetGroupKey()
	group := q.Doc.QuorumsByEvent[group_key]
	delete(group, key)
	if len(group) == 0 {
		delete(q.Doc.QuorumsByEvent, group_key)
	}
}
