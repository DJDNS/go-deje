package deje

import (
	"github.com/campadrenalin/go-deje/serial"
	"github.com/campadrenalin/go-deje/util"
)

type Event struct {
	ParentHash  string                 `json:"parent"`
	HandlerName string                 `json:"handler"`
	Arguments   map[string]interface{} `json:"args"`
}

type EventSet map[string]Event

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

func (e *Event) SetParent(p Event) {
	e.ParentHash = p.Hash()
}

// Get the hash of the Event object.
func (e Event) Hash() string {
	hash, _ := util.HashObject(e)
	return hash
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
func (s EventSet) GetCommonAncestor(A, B *Event) *Event {
	AncestorsA := make(EventSet)
	AncestorsB := make(EventSet)

	AncestorsA.Register(*A)
	AncestorsB.Register(*B)

	for {
		anc, ok := s[A.ParentHash]
		if !ok {
			return nil
		}
		A = &anc
		if AncestorsB.Contains(*A) {
			return A
		} else {
			AncestorsA.Register(*A)
		}

		A, B = B, A
		AncestorsA, AncestorsB = AncestorsB, AncestorsA
	}
}

func (s EventSet) Register(event Event) {
	hash := event.Hash()
	s[hash] = event
}

func (s EventSet) Contains(event Event) bool {
	hash := event.Hash()
	_, ok := s[hash]

	return ok
}

func (s EventSet) GetRoot(tip Event) (event Event, ok bool) {
	event = tip
	ok = true
	for event.ParentHash != "" {
		event, ok = s[event.ParentHash]
		if !ok {
			return
		}
	}
	return
}