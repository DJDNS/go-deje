package manager

import (
	"github.com/campadrenalin/go-deje/model"
	"testing"
)

func TestEventManager_Register(t *testing.T) {
	m := NewEventManager()
	e := model.NewEvent("handler name")

	if m.Contains(e) {
		t.Fatal("m should not contain e yet")
	}
	m.Register(e)
	if !m.Contains(e) {
		t.Fatal("m should contain e")
	}
}

func TestEventManager_Unregister(t *testing.T) {
	m := NewEventManager()
	e := model.NewEvent("handler name")
	m.Register(e)

	if !m.Contains(e) {
		t.Fatal("m should contain e")
	}
	m.Unregister(e)
	if m.Contains(e) {
		t.Fatal("m should not contain e anymore")
	}

	// Should be idempotent
	m.Unregister(e)
	if m.Contains(e) {
		t.Fatal("m should not contain e anymore")
	}
}

func TestEventManager_DeserializeFrom(t *testing.T) {
	m := NewEventManager()
	serial := make(map[string]model.Event)

	// Use fake keys - should be ignored
	serial["energy drink"] = model.NewEvent("on the ground")
	serial["hot dog"] = model.NewEvent("hand outs")
	serial["cell phone"] = model.NewEvent("not my dad")

	m.DeserializeFrom(serial)

	items := m.GetItems()
	for key, value := range serial {
		_, ok := items[key]
		if ok {
			t.Fatalf("Key %s should not have been present in items", key)
		}

		if !m.Contains(value) {
			t.Fatalf("m should contain %#v", value)
		}
	}
}

func TestEventManager_SerializeTo(t *testing.T) {
	m := NewEventManager()
	serial := make(map[string]model.Event)

	events := []model.Event{
		model.NewEvent("A"),
		model.NewEvent("B"),
		model.NewEvent("C"),
	}
	for _, value := range events {
		m.Register(value)
	}

	m.SerializeTo(serial)

	if len(serial) != len(events) {
		t.Fatalf("Expected %d events, got %d", len(events), len(serial))
	}
	for _, value := range events {
		key := value.GetKey()
		serial_event, ok := serial[key]
		if !ok {
			t.Fatalf("Item %s not present", key)
		}
		if !value.Eq(serial_event) {
			t.Fatalf("%#v != %#v", value, serial_event)
		}
	}
}
