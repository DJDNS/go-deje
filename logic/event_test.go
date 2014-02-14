package logic

import "testing"

func TestEvent_GetCommonAncestor_CommonAncestorExists(t *testing.T) {
	d := NewDocument()
	ev_root := d.NewEvent("root")
	ev_childA := d.NewEvent("childA")
	ev_childB := d.NewEvent("childB")

	ev_childA.SetParent(ev_root)
	ev_childB.SetParent(ev_root)

	events := []Event{ev_root, ev_childA, ev_childB}
	for _, ev := range events {
		d.Events.Register(ev.Event)
	}

	anc_ab, err := ev_childA.GetCommonAncestor(ev_childB)
	if err != nil {
		t.Fatal(err)
	}
	if !anc_ab.Eq(ev_root) {
		t.Fatal("Common ancestor of A and B should be root")
	}

	anc_ba, err := ev_childB.GetCommonAncestor(ev_childA)
	if err != nil {
		t.Fatal(err)
	}
	if !anc_ba.Eq(ev_root) {
		t.Fatal("Common ancestor of A and B should be root")
	}
}

func TestEvent_GetCommonAncestor_MissingParent(t *testing.T) {
	d := NewDocument()

	ev_A := d.NewEvent("handler_name")
	ev_A.ParentHash = "blah blah blah"
	ev_B := d.NewEvent("handler_name")

	events := []Event{ev_A, ev_B}
	for _, ev := range events {
		d.Events.Register(ev.Event)
	}

	_, err := ev_A.GetCommonAncestor(ev_B)
	if err == nil {
		t.Fatal("GetCommonAncestor should have failed")
	}
}

func TestEvent_GetCommonAncestor_RootVSFarChild(t *testing.T) {
	d := NewDocument()
	ev_root := d.NewEvent("root")
	ev_childA := d.NewEvent("childA")
	ev_childB := d.NewEvent("childB")

	ev_childA.SetParent(ev_root)
	ev_childB.SetParent(ev_childA)

	events := []Event{ev_root, ev_childA, ev_childB}
	for _, ev := range events {
		d.Events.Register(ev.Event)
	}

	anc_rb, err := ev_root.GetCommonAncestor(ev_childB)
	if err != nil {
		t.Fatal(err)
	}
	if !anc_rb.Eq(ev_root) {
		t.Fatal("Common ancestor of root and B should be root")
	}

	anc_br, err := ev_childB.GetCommonAncestor(ev_root)
	if err != nil {
		t.Fatal(err)
	}
	if !anc_br.Eq(ev_root) {
		t.Fatal("Common ancestor of root and B should be root")
	}
}

func TestEvent_GetCommonAncestor_NoCA(t *testing.T) {
	d := NewDocument()

	ev_A := d.NewEvent("A")
	ev_B := d.NewEvent("B")

	events := []Event{ev_A, ev_B}
	for _, ev := range events {
		d.Events.Register(ev.Event)
	}

	_, err := ev_A.GetCommonAncestor(ev_B)
	if err == nil {
		t.Fatal("GetCommonAncestor should have failed")
	}
}

func TestEvent_GetCommonAncestor_ComparedToSelf(t *testing.T) {
	d := NewDocument()
	ev := d.NewEvent("ev")

	d.Events.Register(ev.Event)

	anc, err := ev.GetCommonAncestor(ev)
	if err != nil {
		t.Fatal(err)
	}
	if !anc.Eq(ev) {
		t.Fatal("Common ancestor of self should be self")
	}
}

func TestEvent_GetRoot_NoElements(t *testing.T) {
	d := NewDocument()
	ev := d.NewEvent("handler_name")
	ev.ParentHash = "blah blah blah" // Not already root
	_, ok := ev.GetRoot()
	if ok {
		t.Fatal("GetRoot should have failed, but returned ok == true")
	}
}

func TestEvent_GetRoot(t *testing.T) {
	d := NewDocument()
	first := d.NewEvent("first")
	second := d.NewEvent("second")
	third := d.NewEvent("third")

	second.SetParent(first)
	third.SetParent(second)

	events := []Event{first, second, third}
	for _, ev := range events {
		d.Events.Register(ev.Event)
	}

	for _, ev := range events {
		found, ok := ev.GetRoot()
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
	first := d.NewEvent("first")
	second := d.NewEvent("second")
	third := d.NewEvent("third")
	fork := d.NewEvent("fork")

	second.SetParent(first)
	third.SetParent(second)
	fork.SetParent(first)

	events := []Event{first, second, third, fork}
	for _, ev := range events {
		d.Events.Register(ev.Event)
	}

	children := first.GetChildren()
	if len(children) != 2 {
		t.Fatal("first has wrong number of children")
	}
	if !(children.Contains(second) && children.Contains(fork)) {
		t.Fatal("first has wrong children", children)
	}

	children = second.GetChildren()
	if len(children) != 1 {
		t.Fatal("second has wrong number of children")
	}
	if !children.Contains(third) {
		t.Fatal("second has wrong children", children)
	}

	children = third.GetChildren()
	if len(children) != 0 {
		t.Fatal("third has wrong number of children")
	}
}
