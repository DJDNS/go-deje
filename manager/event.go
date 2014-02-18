package manager

import "github.com/campadrenalin/go-deje/model"

type EventManager struct {
	genericManager
}

func NewEventManager() *EventManager {
	om := newGenericManager()
	return &EventManager{om}
}

func (em *EventManager) Register(event model.Event) {
	em.register(event)
}

func (em *EventManager) Unregister(event model.Event) {
	em.unregister(event)
}

func (qm *EventManager) DeserializeFrom(items map[string]model.Event) {
	for _, value := range items {
		qm.Register(value)
	}
}

func (qm *EventManager) SerializeTo(items map[string]model.Event) {
	for key, value := range qm.GetItems() {
		ev := value.(model.Event)
		items[key] = ev
	}
}
