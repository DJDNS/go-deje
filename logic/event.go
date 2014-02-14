package logic

import (
	"errors"
	"github.com/campadrenalin/go-deje/model"
)

type Event struct {
	model.Event
	Doc *Document
}

func (doc *Document) NewEvent(handler_name string) Event {
	return Event{
		model.NewEvent(handler_name),
		doc,
	}
}

func (e *Event) SetParent(p Event) {
	e.ParentHash = p.Hash()
}

func (e *Event) Register() {
	e.Doc.Events.Register(e.Event)
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
func (A Event) GetCommonAncestor(B Event) (Event, error) {
	d := A.Doc
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
			if current.Event.ParentHash == "" {
				continue
			}
			parent, ok := d.Events.GetByKey(current.ParentHash)
			if !ok {
				return current, errors.New("Bad parent hash")
			}
			trails <- Event{parent.(model.Event), d}
		default:
			return current, errors.New("No common ancestor")
		}
	}
}

func (tip Event) GetRoot() (event Event, ok bool) {
	var parent model.Manageable
	event = tip
	ok = true
	d := tip.Doc

	for event.Event.ParentHash != "" {
		parent, ok = d.Events.GetByKey(event.Event.ParentHash)
		if !ok {
			return
		}
		event = Event{parent.(model.Event), d}
	}
	return
}

// Get a list of the children of an Event.
func (e Event) GetChildren() model.ManageableSet {
	group_key := e.Hash()
	return e.Doc.Events.GetGroup(group_key)
}
