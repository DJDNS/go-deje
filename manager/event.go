package manager

import "github.com/campadrenalin/go-deje/model"

type EventManager struct {
	GenericManager
}

func NewEventManager() EventManager {
	om := NewGenericManager()
	return EventManager{om}
}

func (em *EventManager) Register(event model.Event) {
	em.register(event)
}

func (em *EventManager) Unregister(event model.Event) {
	em.unregister(event)
}
