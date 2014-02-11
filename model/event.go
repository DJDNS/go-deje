package model

import (
	"errors"
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
// an error.
func (A Event) GetCommonAncestor(d Document, B Event) (Event, error) {
	ancestors := make(map[string]bool)
	trails := make(chan Event, 2)
	var current Event

	trails <- A
	trails <- B

	for {
		select {
		case current = <-trails:
			// Check current
			current_key := current.GetKey()
			if ancestors[current_key] {
				return current, nil
			} else {
				ancestors[current_key] = true
			}

			// Get parent, add to trails
			if current.ParentHash == "" {
				continue
			}
			parent, ok := d.Events.GetByKey(current.ParentHash)
			if !ok {
				return current, errors.New("Bad parent hash")
			}
			trails <- parent.(Event)
		default:
			return current, errors.New("No common ancestor")
		}
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

// Get a list of the children of an Event.
func (e Event) GetChildren(d Document) ManageableSet {
	group_key := e.Hash()
	return d.Events.GetGroup(group_key)
}

// Serialization

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
