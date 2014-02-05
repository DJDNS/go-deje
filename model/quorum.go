package model

import (
	"github.com/campadrenalin/go-deje/serial"
	"github.com/campadrenalin/go-deje/util"
)

// Represents a complete set of approvals for an event.
// Quorums act as bridges between events and timestamps,
// indicating that an event was both common knowledge and
// considered a valid event chain (among others) at one
// time (the timestamp provides the time information).
type Quorum struct {
	EventHash  string
	Signatures map[string]string
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

// Serialization

func QuorumFromSerial(sq serial.Quorum) Quorum {
	return Quorum{
		EventHash:  sq.EventHash,
		Signatures: sq.Signatures,
	}
}

func (q Quorum) ToSerial() serial.Quorum {
	return serial.Quorum{
		EventHash:  q.EventHash,
		Signatures: q.Signatures,
	}
}

func QuorumSetFromObjectManager(om ObjectManager) serial.QuorumSet {
	qs := make(serial.QuorumSet)

	for key, value := range om.GetItems() {
		q, ok := value.(Quorum)
		if ok {
			qs[key] = q.ToSerial()
		}
	}

	return qs
}

func ObjectManagerFromQuorumSet(qs serial.QuorumSet) ObjectManager {
	om := NewObjectManager()

	for _, value := range qs {
		om.Register(QuorumFromSerial(value))
	}

	return om
}
