package deje

import "testing"

func TestEventSet_GetRoot_NoElements(t *testing.T) {
	set := make(EventSet)
	ev := NewEvent("handler_name")
	ev.ParentHash = "blah blah blah" // Not already root

	_, ok := set.GetRoot(ev)
	if ok {
		t.Fatal("GetRoot should have failed, but returned ok == true")
	}
}

func TestEventSet_GetRoot(t *testing.T) {
	set := make(EventSet)
	first := NewEvent("first")
	second := NewEvent("second")
	third := NewEvent("third")

	second.SetParent(first)
	third.SetParent(second)

	events := []Event{first, second, third}
	for _, ev := range events {
		set.Register(ev)
	}

	for _, ev := range events {
		found, ok := set.GetRoot(ev)
		if !ok {
			t.Fatal("GetRoot failed")
		}
		if found.HandlerName != "first" {
			t.Fatal("Did not get correct event")
		}
	}
}

func TestEventSetContains(t *testing.T) {
	set := make(EventSet)
	first := NewEvent("first")
	second := NewEvent("second")
	third := NewEvent("third")

	second.SetParent(first)
	third.SetParent(second)

	events := []Event{first, third} // Every event but second
	for _, ev := range events {
		set.Register(ev)
	}

    if ! set.Contains(first) {
        t.Fatal("set should contain first")
    }
    if ! set.Contains(third) {
        t.Fatal("set should contain third")
    }

    if set.Contains(second) {
        t.Fatal("set should contain second")
    }
}
