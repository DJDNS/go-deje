package document

import "github.com/campadrenalin/go-deje/util"

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
