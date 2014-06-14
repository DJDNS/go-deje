package logic

import (
	"errors"
	"github.com/campadrenalin/go-deje/model"
	"github.com/campadrenalin/go-deje/state"
)

type Event struct {
	model.Event
	Doc *Document
}
type EventSet map[string]*Event

func (es EventSet) Contains(ev Event) bool {
	_, ok := es[ev.GetKey()]
	return ok
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

// Register with the Doc. This stores it in a hash-based location,
// so do not make changes to an Event after it has been registered.
func (e *Event) Register() {
	key := e.GetKey()
	e.Doc.Events[key] = e

	group_key := e.GetGroupKey()
	group, ok := e.Doc.EventsByParent[group_key]
	if !ok {
		group = make(EventSet)
		e.Doc.EventsByParent[group_key] = group
	}
	group[key] = e
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
			parent, ok := d.Events[current.ParentHash]
			if !ok {
				return current, errors.New("Bad parent hash")
			}
			trails <- Event{parent.Event, d}
		default:
			return current, errors.New("No common ancestor")
		}
	}
}

// Returns whether two events are on compatible forks.
//
// This means that they are on the same fork. This means
// one is the parent of the other (or the events are equal).
func (A Event) CompatibleWith(B Event) (bool, error) {
	parent, err := A.GetCommonAncestor(B)
	if err != nil {
		return false, err
	}

	return (parent.Eq(A.Event) || parent.Eq(B.Event)), nil
}

// Traverse up the chain of parents until there's no more to traverse.
func (tip *Event) GetRoot() (event *Event, ok bool) {
	var parent *Event
	event = tip
	ok = true
	d := tip.Doc

	for event.Event.ParentHash != "" {
		parent, ok = d.Events[event.Event.ParentHash]
		if !ok {
			return
		}
		event = parent
	}
	return
}

// Get a list of the children of an Event.
func (e Event) GetChildren() EventSet {
	group_key := e.GetKey()
	return e.Doc.EventsByParent[group_key]
}

// Translate this event into a set of primitives.
//
// Custom events may fail to translate, due to failures
// in the Lua interpreter environment, or attempting to get
// the event's primitives when the document's state is not
// on the event's parent.
//
// Builtin events (SET and DELETE) should always be able to
// be translated into a primitive, regardless of doc state,
// as long as the event's properties are sufficient to populate
// the struct primitive.
func (e Event) getPrimitives() ([]state.Primitive, error) {
	path_interface, ok := e.Arguments["path"]
	if !ok {
		return nil, errors.New("No path provided")
	}
	path, ok := path_interface.([]interface{})
	if !ok {
		return nil, errors.New("Bad path value")
	}
	value, ok := e.Arguments["value"]
	if !ok {
		return nil, errors.New("No value provided")
	}
	primitives := []state.Primitive{
		&state.SetPrimitive{
			Path:  path,
			Value: value,
		},
	}
	return primitives, nil
}

// Attempt to apply this event to the current document state.
//
// Does not check that the document is at the Event's parent
// before attempting to apply primitives.
func (e Event) Apply() error {
	primitives, err := e.getPrimitives()
	if err != nil {
		return err
	}
	for _, primitive := range primitives {
		err = e.Doc.State.Apply(primitive)
		if err != nil {
			return err
		}
	}
	return nil
}

// Attempt to navigate the DocumentState to this Event.
//
// Somewhat analogous to git checkout.
func (e Event) Goto() error {
	d := e.Doc
	d.State.Reset()
	if e.Event.ParentHash != "" {
		parent, ok := d.Events[e.Event.ParentHash]
		if !ok {
			return errors.New("Could not get parent")
		}
		logic_parent := Event{parent.Event, d}
		err := logic_parent.Goto()
		if err != nil {
			return err
		}
	}
	return e.Apply()
}
