package manager

import "github.com/campadrenalin/go-deje/model"

// Manager that stores Events.
//
// Events are grouped by their parent hash, and keyed on their own hash.
type EventManager struct {
	genericManager
}

func NewEventManager() *EventManager {
	om := newGenericManager()
	return &EventManager{om}
}

// Register an Event. Afterwards, it can be retrieved by key or group.
func (em *EventManager) Register(event model.Event) {
	em.register(event)
}

// Remove an Event from the internal registries.
func (em *EventManager) Unregister(event model.Event) {
	em.unregister(event)
}

// Import a set of Events, registering all of them.
func (qm *EventManager) DeserializeFrom(items map[string]model.Event) {
	for _, value := range items {
		qm.Register(value)
	}
}

// Export a set of Events to a raw map.
func (qm *EventManager) SerializeTo(items map[string]model.Event) {
	for key, value := range qm.GetItems() {
		ev := value.(model.Event)
		items[key] = ev
	}
}
