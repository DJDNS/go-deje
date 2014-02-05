package deje

import (
	"github.com/campadrenalin/go-deje/serial"
	"github.com/campadrenalin/go-deje/util"
)

type Event struct {
	ParentHash  string
	HandlerName string
	Arguments   map[string]interface{}
}

func NewEvent(hname string) Event {
	return Event{
		ParentHash:  "",
		HandlerName: hname,
		Arguments:   make(map[string]interface{}),
	}
}

func EventFromSerial(se serial.Event) Event {
	return Event{
		ParentHash:  se.ParentHash,
		HandlerName: se.HandlerName,
		Arguments:   se.Arguments,
	}
}

func (e *Event) ToSerial() serial.Event {
	return serial.Event{
		ParentHash:  e.ParentHash,
		HandlerName: e.HandlerName,
		Arguments:   e.Arguments,
	}
}

func (e Event) GetKey() string {
	return e.Hash()
}
func (e Event) GetGroupKey() string {
	return e.ParentHash
}
func (e Event) Eq(other Manageable) bool {
	other_event, ok := other.(Event)
	if !ok {
		return false
	}
	return e.Hash() == other_event.Hash()
}

// Get the hash of the Event object.
func (e Event) Hash() string {
	hash, _ := util.HashObject(e)
	return hash
}

func (e *Event) SetParent(p Event) {
	e.ParentHash = p.Hash()
}

// Given a set of Events, and two specific ones to trace,
// find the most recent common parent between the two chains.
//
// If A is the common ancestor, it is a parent of B. And vice
// versa - if B is the common ancestor, it is a parent of A.
// There is also the corner case where you compare an Event
// against itself, and get that same Event. However, if the
// common ancestor is neither A nor B, than the two Events
// are not in the same chain of history, and must be considered
// incompatible branches of the Event chain.
//
// There may not be a common ancestor. In this event, we return
// a nil pointer.
func (A Event) GetCommonAncestor(d Document, B Event) (Event, bool) {
	AncestorsA := make(map[string]bool)
	AncestorsB := make(map[string]bool)

	AncestorsA[A.Hash()] = true
	AncestorsB[B.Hash()] = true

	for {
		anc, ok := d.Events.GetByKey(A.ParentHash)
		if !ok {
			return Event{}, false
		}
		A = anc.(Event)
		if AncestorsB[A.Hash()] {
			return A, true
		} else {
			AncestorsA[A.Hash()] = true
		}

		A, B = B, A
		AncestorsA, AncestorsB = AncestorsB, AncestorsA
	}
}

func (tip Event) GetRoot(d Document) (event Event, ok bool) {
	var parent Manageable
	event = tip
	ok = true
	for event.ParentHash != "" {
		parent, ok = d.Events.GetByKey(event.ParentHash)
		if !ok {
			return
		}
		event = parent.(Event)
	}
	return
}
