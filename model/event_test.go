package model

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

func TestEvent_GetChildren(t *testing.T) {
	d := NewDocument()
	first := NewEvent("first")
	second := NewEvent("second")
	third := NewEvent("third")
	fork := NewEvent("fork")

	second.SetParent(first)
	third.SetParent(second)
	fork.SetParent(first)

	events := []Event{first, second, third, fork}
	for _, ev := range events {
		d.Events.Register(ev)
	}

	children := first.GetChildren(d)
	if len(children) != 2 {
		t.Fatal("first has wrong number of children")
	}
	if !(children.Contains(second) && children.Contains(fork)) {
		t.Fatal("first has wrong children", children)
	}

	children = second.GetChildren(d)
	if len(children) != 1 {
		t.Fatal("second has wrong number of children")
	}
	if !children.Contains(third) {
		t.Fatal("second has wrong children", children)
	}

	children = third.GetChildren(d)
	if len(children) != 0 {
		t.Fatal("third has wrong number of children")
	}
}
