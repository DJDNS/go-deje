package model

import "github.com/campadrenalin/go-deje/serial"

type EventManager struct {
	ObjectManager
}

func NewEventManager() EventManager {
	om := NewObjectManager()
	return EventManager{om}
}

func (em *EventManager) Register(event Event) {
	em.register(event)
}

func (em *EventManager) Unregister(event Event) {
	em.unregister(event)
}

func (em *EventManager) DeserializeFrom(items map[string]serial.Event) {
	for _, value := range items {
		em.Register(EventFromSerial(value))
	}
}

func (em *EventManager) SerializeTo(items map[string]serial.Event) {
	for key, value := range em.GetItems() {
		ev := value.(Event)
		items[key] = ev.ToSerial()
	}
}
