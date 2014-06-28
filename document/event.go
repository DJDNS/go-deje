package document

import (
	"errors"

	"github.com/campadrenalin/go-deje/state"
	"github.com/campadrenalin/go-deje/util"
)

// An Event is an action that can be applied to a DEJE doc,
// including a set of parameters. In practice, custom Event
// types may be defined for a document, as well as permissions
// for which users are allowed to perform which types of Events.
type Event struct {
	Doc         *Document              `json:"-"`
	ParentHash  string                 `json:"parent"`
	HandlerName string                 `json:"handler"`
	Arguments   map[string]interface{} `json:"args"`
}

type EventSet map[string]*Event

func (es EventSet) Contains(ev Event) bool {
	_, ok := es[ev.GetKey()]
	return ok
}

func NewEvent(hname string) Event {
	return Event{
		ParentHash:  "",
		HandlerName: hname,
		Arguments:   make(map[string]interface{}),
	}
}

func (doc *Document) NewEvent(handler_name string) Event {
	ev := NewEvent(handler_name)
	ev.Doc = doc
	return ev
}

func (e Event) GetKey() string {
	return e.Hash()
}
func (e Event) GetGroupKey() string {
	return e.ParentHash
}
func (e Event) Eq(other Event) bool {
	return e.Hash() == other.Hash()
}

// Get the hash of the Event object.
func (e Event) Hash() string {
	hash, _ := util.HashObject(e)
	return hash
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

// Unregister from the Doc. This also cleans up empty groups.
func (e *Event) Unregister() {
	key := e.GetKey()
	delete(e.Doc.Events, key)

	group_key := e.GetGroupKey()
	group := e.Doc.EventsByParent[group_key]
	delete(group, key)
	if len(group) == 0 {
		delete(e.Doc.EventsByParent, group_key)
	}
}

// Convenience function. Panics if e.Doc is nil.
func (e *Event) GetParent() (*Event, bool) {
	p, ok := e.Doc.Events[e.ParentHash]
	return p, ok
}

// Get a linear history of the Events leading up to the given Event.
//
// If we fail to find a parent at any point, we return (nil, false).
func (e *Event) GetHistory() ([]*Event, bool) {
	history := []*Event{e}
	current := e
	for {
		if current.ParentHash == "" {
			return history, true
		}
		parent, ok := current.GetParent()
		if ok {
			// Prepend parent onto history
			history = append([]*Event{parent}, history...)
			current = parent
		} else {
			return nil, false
		}
	}
	return history, true
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
func (A *Event) GetCommonAncestor(B *Event) (*Event, error) {
	d := A.Doc
	ancestors := make(map[string]bool)
	trails := make(chan *Event, 2)
	var current *Event

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
			parent, ok := d.Events[current.ParentHash]
			if !ok {
				return current, errors.New("Bad parent hash")
			}
			trails <- parent
		default:
			return current, errors.New("No common ancestor")
		}
	}
}

// Returns whether two events are on compatible forks.
//
// This means that they are on the same fork. This means
// one is the parent of the other (or the events are equal).
func (A *Event) CompatibleWith(B *Event) (bool, error) {
	parent, err := A.GetCommonAncestor(B)
	if err != nil {
		return false, err
	}

	return (parent.Eq(*A) || parent.Eq(*B)), nil
}

// Traverse up the chain of parents until there's no more to traverse.
func (tip *Event) GetRoot() (event *Event, ok bool) {
	var parent *Event
	event = tip
	ok = true
	d := tip.Doc

	for event.ParentHash != "" {
		parent, ok = d.Events[event.ParentHash]
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
	if e.ParentHash != "" {
		parent, ok := d.Events[e.ParentHash]
		if !ok {
			return errors.New("Could not get parent")
		}
		err := parent.Goto()
		if err != nil {
			return err
		}
	}
	return e.Apply()
}
