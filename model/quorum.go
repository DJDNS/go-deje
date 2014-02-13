package model

import "github.com/campadrenalin/go-deje/util"

// Represents a complete set of approvals for an event.
// Quorums act as bridges between events and timestamps,
// indicating that an event was both common knowledge and
// considered a valid event chain (among others) at one
// time (the timestamp provides the time information).
type Quorum struct {
	EventHash  string            `json:"event_hash"`
	Signatures map[string]string `json:"sigs"`
}

type QuorumSet map[string]Quorum

func NewQuorum(evhash string) Quorum {
	return Quorum{
		EventHash:  evhash,
		Signatures: make(map[string]string),
	}
}

func (q Quorum) GetKey() string {
	return q.Hash()
}
func (q Quorum) GetGroupKey() string {
	return q.EventHash
}
func (q Quorum) Eq(other Manageable) bool {
	other_quorum, ok := other.(Quorum)
	if !ok {
		return false
	}
	return q.Hash() == other_quorum.Hash()
}

// Get the hash of the Quorum object.
func (q Quorum) Hash() string {
	hash, _ := util.HashObject(q)
	return hash
}
