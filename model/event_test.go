package deje

import "testing"

func TestEventSet_GetRoot_NoElements(t *testing.T) {
	d := NewDocument()
	ev := NewEvent("handler_name")
	ev.ParentHash = "blah blah blah" // Not already root

	_, ok := ev.GetRoot(d)
	if ok {
		t.Fatal("GetRoot should have failed, but returned ok == true")
	}
}

func TestEvent_GetRoot(t *testing.T) {
	d := NewDocument()
	first := NewEvent("first")
	second := NewEvent("second")
	third := NewEvent("third")

	second.SetParent(first)
	third.SetParent(second)

	events := []Event{first, second, third}
	for _, ev := range events {
		d.Events.Register(ev)
	}

	for _, ev := range events {
		found, ok := ev.GetRoot(d)
		if !ok {
			t.Fatal("GetRoot failed")
		}
		if found.HandlerName != "first" {
			t.Fatal("Did not get correct event")
		}
	}
}
