package model

import "github.com/campadrenalin/go-deje/util"

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
