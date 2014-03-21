package logic

import (
	"github.com/campadrenalin/go-deje/state"
	"reflect"
	"testing"
)

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
	if !anc_ab.Eq(ev_root.Event) {
		t.Fatal("Common ancestor of A and B should be root")
	}

	anc_ba, err := ev_childB.GetCommonAncestor(ev_childA)
	if err != nil {
		t.Fatal(err)
	}
	if !anc_ba.Eq(ev_root.Event) {
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
	if !anc_rb.Eq(ev_root.Event) {
		t.Fatal("Common ancestor of root and B should be root")
	}

	anc_br, err := ev_childB.GetCommonAncestor(ev_root)
	if err != nil {
		t.Fatal(err)
	}
	if !anc_br.Eq(ev_root.Event) {
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
	if !anc.Eq(ev.Event) {
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

func assert_ev_compatible(A, B Event, t *testing.T) {
	compatible, err := A.CompatibleWith(B)
	if err != nil {
		t.Fatal(err)
	}
	if !compatible {
		t.Fatalf("%s should be compatible with %s", A.Event.HandlerName, B.Event.HandlerName)
	}
}

func assert_ev_incompatible(A, B Event, t *testing.T) {
	compatible, err := A.CompatibleWith(B)
	if err != nil {
		t.Fatal(err)
	}
	if compatible {
		t.Fatalf("%s should be incompatible with %s", A.Event.HandlerName, B.Event.HandlerName)
	}
}

func TestEvent_CompatibleWith(t *testing.T) {
	d := NewDocument()
	first := d.NewEvent("first")
	second := d.NewEvent("second")
	third := d.NewEvent("third")
	fork := d.NewEvent("fork")

	second.SetParent(first)
	third.SetParent(second)
	fork.SetParent(first)

	_, err := fork.CompatibleWith(first)
	if err == nil {
		t.Fatal("Events not registered yet, should have failed")
	}

	for _, e := range []Event{first, second, third, fork} {
		e.Register()
	}

	assert_ev_compatible(first, second, t)
	assert_ev_compatible(first, third, t)
	assert_ev_compatible(first, fork, t)

	assert_ev_incompatible(fork, second, t)
	assert_ev_incompatible(fork, third, t)
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
	if !(children.Contains(second.Event) && children.Contains(fork.Event)) {
		t.Fatal("first has wrong children", children)
	}

	children = second.GetChildren()
	if len(children) != 1 {
		t.Fatal("second has wrong number of children")
	}
	if !children.Contains(third.Event) {
		t.Fatal("second has wrong children", children)
	}

	children = third.GetChildren()
	if len(children) != 0 {
		t.Fatal("third has wrong number of children")
	}
}

type eventToPrimitivesTest struct {
	HandlerName string
	Arguments   map[string]interface{}
	Expected    []state.Primitive
	ShouldFail  bool
	FailureMsg  string
}

func (test *eventToPrimitivesTest) Run(t *testing.T) {
	d := NewDocument()
	setter := d.NewEvent(test.HandlerName)
	setter.Arguments = test.Arguments

	primitives, err := setter.getPrimitives()
	if test.ShouldFail && err == nil {
		t.Error(test.FailureMsg)
		t.Fatal("Event.getPrimitives should have failed, didn't")
	} else if !test.ShouldFail && err != nil {
		t.Error(test.FailureMsg)
		t.Fatal(err)
	}

	if !reflect.DeepEqual(primitives, test.Expected) {
		t.Error(test.FailureMsg)
		t.Fatalf("Expected %#v, got %#v", test.Expected, primitives)
	}
}

func TestEvent_getPrimitives_Set(t *testing.T) {
	tests := []eventToPrimitivesTest{
		eventToPrimitivesTest{
			"SET",
			map[string]interface{}{
				"path":  []interface{}{"hello"},
				"value": "world",
			},
			[]state.Primitive{
				&state.SetPrimitive{
					Path:  []interface{}{"hello"},
					Value: "world",
				},
			},
			false,
			"Basic SET event with reasonable params",
		},
		eventToPrimitivesTest{
			"SET",
			map[string]interface{}{
				"value": "world",
			},
			nil, true,
			"SET with no path",
		},
		eventToPrimitivesTest{
			"SET",
			map[string]interface{}{
				"path":  7,
				"value": "world",
			},
			nil, true,
			"SET with bad path",
		},
		eventToPrimitivesTest{
			"SET",
			map[string]interface{}{
				"path": []interface{}{"hello"},
			},
			nil, true,
			"SET with no value", // Like an obscure altcoin
		},
	}
	for _, test := range tests {
		test.Run(t)
	}
}

func TestEvent_Apply(t *testing.T) {
	d := NewDocument()
	// TODO: Test doomed-to-failure/invalid events
	ev := d.NewEvent("SET")
	ev.Arguments["path"] = []interface{}{}
	ev.Arguments["value"] = map[string]interface{}{
		"hello": "world",
	}
	primitives, err := ev.getPrimitives()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Primitives: %#v", primitives)

	err = ev.Apply()
	expected_export := ev.Arguments["value"]
	exported := d.State.Export()
	if !reflect.DeepEqual(exported, expected_export) {
		t.Fatalf("Expected %#v, got %#v", expected_export, exported)
	}
}

func TestEvent_Goto(t *testing.T) {
	d := NewDocument()
	// TODO: Test doomed-to-failure/invalid events
	ev_root := d.NewEvent("SET")
	ev_root.Arguments["path"] = []interface{}{}
	ev_root.Arguments["value"] = map[string]interface{}{
		"hello": map[string]interface{}{},
	}

	ev_child := d.NewEvent("SET")
	ev_child.Arguments["path"] = []interface{}{"hello", "to the"}
	ev_child.Arguments["value"] = "world"
	ev_child.SetParent(ev_root)

	ev_fork := d.NewEvent("SET")
	ev_fork.Arguments["path"] = []interface{}{"hello", "little"}
	ev_fork.Arguments["value"] = "fork"
	ev_fork.SetParent(ev_root)

	ev_root.Register()
	ev_child.Register()
	ev_fork.Register()

	// Test first
	err := ev_child.Goto()
	if err != nil {
		t.Fatal(err)
	}
	expected_export := map[string]interface{}{
		"hello": map[string]interface{}{
			"to the": "world",
		},
	}
	exported := d.State.Export()
	if !reflect.DeepEqual(exported, expected_export) {
		t.Fatalf("Expected %#v, got %#v", expected_export, exported)
	}

	// Test switch
	err = ev_fork.Goto()
	if err != nil {
		t.Fatal(err)
	}
	expected_export = map[string]interface{}{
		"hello": map[string]interface{}{
			"little": "fork",
		},
	}
	exported = d.State.Export()
	if !reflect.DeepEqual(exported, expected_export) {
		t.Fatalf("Expected %#v, got %#v", expected_export, exported)
	}
}
