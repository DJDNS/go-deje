package document

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/DJDNS/go-deje/state"
	"github.com/stretchr/testify/assert"
)

func TestEvent_Serialize(t *testing.T) {
	ev := NewEvent("handler_name")
	ev.Arguments["hello"] = []interface{}{"world", 5}
	ev.Arguments["before"] = nil

	serialized, err := json.Marshal(ev)
	if err != nil {
		t.Fatal("Serialization failed")
	}
	expected := []byte("{" +
		"\"parent\":\"\"," +
		"\"handler\":\"handler_name\"," +
		"\"args\":{" +
		"\"before\":null," +
		"\"hello\":[\"world\",5]" +
		"}" +
		"}")
	if !bytes.Equal(serialized, expected) {
		t.Fatal(string(serialized))
	}
}

func TestEvent_GetKey(t *testing.T) {
	ev := NewEvent("handler_name")
	ev.Arguments["hello"] = []interface{}{"world", 5}
	ev.Arguments["before"] = nil

	key := ev.GetKey()

	// Obtained via:
	// echo -n '{"parent":"","handler":"handler_name","args":{"before":null,"hello":["world",5]}}' | sha1sum
	expected := "86e5db5fcf8c749146f2adcc23c728769ef2bd98"

	if key != expected {
		t.Fatalf("Expected %v, got %v", expected, key)
	}
}

func TestEvent_GetGroupKey(t *testing.T) {
	ev := NewEvent("SET")

	gk := ev.GetGroupKey()
	if gk != "" {
		t.Fatalf("Expected empty group key, got %v", gk)
	}

	expected := "Hurpdeburp"
	ev.ParentHash = expected
	gk = ev.GetGroupKey()
	if gk != expected {
		t.Fatalf("Expected group key %v, got %v", expected, gk)
	}
}

func TestEvent_Eq(t *testing.T) {
	A := NewEvent("hello")
	B := NewEvent("hello")
	C := NewEvent("hello")
	D := NewEvent("hello")

	if !(A.Eq(A) && A.Eq(B) && A.Eq(C) && A.Eq(D)) {
		t.Fatal("Freshly initialized events are not equal")
	}

	B.ParentHash = "Ezekiel Wigglesworth"
	if A.Eq(B) {
		t.Fatal("A should not equal B")
	}

	C.HandlerName = "Ezekiel Wigglesworth"
	if A.Eq(B) {
		t.Fatal("A should not equal C")
	}

	D.Arguments["Ezekiel"] = "Wigglesworth"
	if A.Eq(B) {
		t.Fatal("A should not equal D")
	}
}

func TestEvent_Register(t *testing.T) {
	d := NewDocument()
	ev_root := d.NewEvent("root")
	ev_childA := d.NewEvent("childA")
	ev_childB := d.NewEvent("childB")

	ev_childA.SetParent(ev_root)
	ev_childB.SetParent(ev_root)

	events := []*Event{&ev_root, &ev_childA, &ev_childB}
	for _, ev := range events {
		ev.Register()
	}

	// Test that main set registered correctly
	expected_events := EventSet{
		ev_root.GetKey():   &ev_root,
		ev_childA.GetKey(): &ev_childA,
		ev_childB.GetKey(): &ev_childB,
	}
	if !reflect.DeepEqual(d.Events, expected_events) {
		t.Fatalf("Expected %#v\nGot %#v", expected_events, d.Events)
	}

	// Test that groupings registered correctly
	expected_groups := map[string]EventSet{
		"": EventSet{
			ev_root.GetKey(): &ev_root,
		},
		ev_root.GetKey(): EventSet{
			ev_childA.GetKey(): &ev_childA,
			ev_childB.GetKey(): &ev_childB,
		},
	}
	if !reflect.DeepEqual(d.EventsByParent, expected_groups) {
		t.Fatalf("Expected %#v\nGot %#v", expected_groups, d.EventsByParent)
	}
}

func TestEvent_Unregister(t *testing.T) {
	d := NewDocument()
	ev_root := d.NewEvent("root")
	ev_childA := d.NewEvent("childA")
	ev_childB := d.NewEvent("childB")

	ev_childA.SetParent(ev_root)
	ev_childB.SetParent(ev_root)

	events := []*Event{&ev_root, &ev_childA, &ev_childB}
	for _, ev := range events {
		ev.Register()
	}

	// Unregister childB and check results
	ev_childB.Unregister()
	expected_events := EventSet{
		ev_root.GetKey():   &ev_root,
		ev_childA.GetKey(): &ev_childA,
	}
	if !reflect.DeepEqual(d.Events, expected_events) {
		t.Fatalf("Expected %#v\nGot %#v", expected_events, d.Events)
	}
	expected_groups := map[string]EventSet{
		"": EventSet{
			ev_root.GetKey(): &ev_root,
		},
		ev_root.GetKey(): EventSet{
			ev_childA.GetKey(): &ev_childA,
		},
	}
	if !reflect.DeepEqual(d.EventsByParent, expected_groups) {
		t.Fatalf("Expected %#v\nGot %#v", expected_groups, d.EventsByParent)
	}

	// Unregister such that a group becomes empty
	ev_root.Unregister()
	expected_events = EventSet{
		ev_childA.GetKey(): &ev_childA,
	}
	if !reflect.DeepEqual(d.Events, expected_events) {
		t.Fatalf("Expected %#v\nGot %#v", expected_events, d.Events)
	}
	expected_groups = map[string]EventSet{
		ev_root.GetKey(): EventSet{
			ev_childA.GetKey(): &ev_childA,
		},
	}
	if !reflect.DeepEqual(d.EventsByParent, expected_groups) {
		t.Fatalf("Expected %#v\nGot %#v", expected_groups, d.EventsByParent)
	}

	// Make sure that Unregistering multiple times is okay
	ev_root.Unregister()
	if !reflect.DeepEqual(d.Events, expected_events) {
		t.Fatalf("Expected %#v\nGot %#v", expected_events, d.Events)
	}
	if !reflect.DeepEqual(d.EventsByParent, expected_groups) {
		t.Fatalf("Expected %#v\nGot %#v", expected_groups, d.EventsByParent)
	}
}

func TestEvent_GetHistory(t *testing.T) {
	d := NewDocument()
	var nil_event_slice []*Event
	ev_first := d.NewEvent("1")
	ev_second := d.NewEvent("2")
	ev_third := d.NewEvent("3")
	ev_fork := d.NewEvent("fork")
	ev_no_parent := d.NewEvent("no parent")

	// Linear chain
	ev_second.SetParent(ev_first)
	ev_third.SetParent(ev_second)
	ev_fork.SetParent(ev_first)              // with hanger-on
	ev_no_parent.ParentHash = "fiddlesticks" // Manually screw this up

	events := []*Event{
		&ev_first, &ev_second, &ev_third,
		&ev_fork, &ev_no_parent,
	}
	for _, ev := range events {
		ev.Register()
	}

	expect := func(t *testing.T, e *Event, history []*Event, ok bool) {
		got_history, got_ok := e.GetHistory()
		t.Logf("Event %s", e.HandlerName)
		assert.Equal(t, ok, got_ok)
		assert.Equal(t, history, got_history)
	}
	expect(t, &ev_first, []*Event{&ev_first}, true)
	expect(t, &ev_second, []*Event{&ev_first, &ev_second}, true)
	expect(t, &ev_third, []*Event{&ev_first, &ev_second, &ev_third}, true)
	expect(t, &ev_fork, []*Event{&ev_first, &ev_fork}, true)
	expect(t, &ev_no_parent, nil_event_slice, false)
}

func TestEvent_GetCommonAncestor_CommonAncestorExists(t *testing.T) {
	d := NewDocument()
	ev_root := d.NewEvent("root")
	ev_childA := d.NewEvent("childA")
	ev_childB := d.NewEvent("childB")

	ev_childA.SetParent(ev_root)
	ev_childB.SetParent(ev_root)

	events := []*Event{&ev_root, &ev_childA, &ev_childB}
	for _, ev := range events {
		ev.Register()
	}

	anc_ab, err := ev_childA.GetCommonAncestor(&ev_childB)
	if err != nil {
		t.Fatal(err)
	}
	if !anc_ab.Eq(ev_root) {
		t.Error("Common ancestor of A and B should be root")
		t.Fatalf("Expected %#v, got %#v", ev_root, anc_ab)
	}

	anc_ba, err := ev_childB.GetCommonAncestor(&ev_childA)
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

	events := []*Event{&ev_A, &ev_B}
	for _, ev := range events {
		ev.Register()
	}

	_, err := ev_A.GetCommonAncestor(&ev_B)
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

	events := []*Event{&ev_root, &ev_childA, &ev_childB}
	for _, ev := range events {
		ev.Register()
	}

	anc_rb, err := ev_root.GetCommonAncestor(&ev_childB)
	if err != nil {
		t.Fatal(err)
	}
	if !anc_rb.Eq(ev_root) {
		t.Fatal("Common ancestor of root and B should be root")
	}

	anc_br, err := ev_childB.GetCommonAncestor(&ev_root)
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

	events := []*Event{&ev_A, &ev_B}
	for _, ev := range events {
		ev.Register()
	}

	_, err := ev_A.GetCommonAncestor(&ev_B)
	if err == nil {
		t.Fatal("GetCommonAncestor should have failed")
	}
}

func TestEvent_GetCommonAncestor_ComparedToSelf(t *testing.T) {
	d := NewDocument()
	ev := d.NewEvent("ev")

	ev.Register()

	anc, err := ev.GetCommonAncestor(&ev)
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

func assert_ev_compatible(A, B *Event, t *testing.T) {
	compatible, err := A.CompatibleWith(B)
	if err != nil {
		t.Fatal(err)
	}
	if !compatible {
		t.Fatalf("%s should be compatible with %s", A.HandlerName, B.HandlerName)
	}
}

func assert_ev_incompatible(A, B *Event, t *testing.T) {
	compatible, err := A.CompatibleWith(B)
	if err != nil {
		t.Fatal(err)
	}
	if compatible {
		t.Fatalf("%s should be incompatible with %s", A.HandlerName, B.HandlerName)
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

	_, err := fork.CompatibleWith(&first)
	if err == nil {
		t.Fatal("Events not registered yet, should have failed")
	}

	for _, ev := range []*Event{&first, &second, &third, &fork} {
		ev.Register()
	}

	assert_ev_compatible(&first, &second, t)
	assert_ev_compatible(&first, &third, t)
	assert_ev_compatible(&first, &fork, t)

	assert_ev_incompatible(&fork, &second, t)
	assert_ev_incompatible(&fork, &third, t)
}

func TestEvent_GetRoot(t *testing.T) {
	d := NewDocument()
	first := d.NewEvent("first")
	second := d.NewEvent("second")
	third := d.NewEvent("third")

	second.SetParent(first)
	third.SetParent(second)

	events := []*Event{&first, &second, &third}
	for _, ev := range events {
		ev.Register()
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
		ev.Register()
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

func TestEvent_getPrimitives_Delete(t *testing.T) {
	tests := []eventToPrimitivesTest{
		eventToPrimitivesTest{
			"DELETE",
			map[string]interface{}{
				"path": []interface{}{"hello"},
			},
			[]state.Primitive{
				&state.DeletePrimitive{
					Path: []interface{}{"hello"},
				},
			},
			false,
			"Basic DELETE event with reasonable params",
		},
		eventToPrimitivesTest{
			"DELETE",
			map[string]interface{}{},
			nil, true,
			"DELETE with no path",
		},
		eventToPrimitivesTest{
			"DELETE",
			map[string]interface{}{
				"path": 7,
			},
			nil, true,
			"DELETE with bad path",
		},
	}
	for _, test := range tests {
		test.Run(t)
	}
}

func TestEvent_getPrimitives_Custom(t *testing.T) {
	tests := []eventToPrimitivesTest{
		eventToPrimitivesTest{
			"some custom event",
			map[string]interface{}{},
			nil, true,
			"Custom events aren't supported yet",
		},
	}
	for _, test := range tests {
		test.Run(t)
	}
}

func TestEvent_Apply(t *testing.T) {
	d := NewDocument()
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

func TestEvent_Apply_BadPrimitives(t *testing.T) {
	d := NewDocument()
	ev := d.NewEvent("SET") // No parameters!

	_, err := ev.getPrimitives()
	if err == nil {
		t.Fatal("ev.getPrimitives should have failed!")
	}
	err = ev.Apply()
	if err == nil {
		t.Fatal("ev.Apply should fail when ev.getPrimitives fails")
	}
}

func TestEvent_Apply_UnapplyablePrimitive(t *testing.T) {
	d := NewDocument()
	ev := d.NewEvent("SET")

	// Invalid event path - does not exist in blank state
	ev.Arguments["path"] = []interface{}{"this", "that"}
	ev.Arguments["value"] = "the other thing"

	err := ev.Apply()
	if err == nil {
		t.Fatal("ev.Apply should fail with unapplyable primitives")
	}
}

func TestEvent_Goto(t *testing.T) {
	d := NewDocument()
	ev_root := d.NewEvent("SET")
	ev_root.Arguments["path"] = []interface{}{"know"}
	ev_root.Arguments["value"] = "fashion's a stranger"

	ev_child := d.NewEvent("SET")
	ev_child.Arguments["path"] = []interface{}{"friend"}
	ev_child.Arguments["value"] = "fashion is danger"
	ev_child.SetParent(ev_root)

	ev_fork := d.NewEvent("SET")
	ev_fork.Arguments["path"] = []interface{}{"posing"}
	ev_fork.Arguments["value"] = "a threat"
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
		"know":   "fashion's a stranger",
		"friend": "fashion is danger",
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
		"know":   "fashion's a stranger",
		"posing": "a threat",
	}
	exported = d.State.Export()
	if !reflect.DeepEqual(exported, expected_export) {
		t.Fatalf("Expected %#v, got %#v", expected_export, exported)
	}
}

func TestEvent_Goto_MissingParent(t *testing.T) {
	d := NewDocument()
	ev_root := d.NewEvent("SET")
	ev_child := d.NewEvent("SET")

	// Set up parentage, but only register child
	ev_child.SetParent(ev_root)
	ev_child.Register()

	err := ev_child.Goto()
	if err == nil {
		t.Fatal("Goto with unreachable heritage should fail!")
	}
}

func TestEvent_Goto_BadParent(t *testing.T) {
	d := NewDocument()
	ev_root := d.NewEvent("SET")
	ev_child := d.NewEvent("SET")

	// Invalid root event path
	ev_root.Arguments["path"] = []interface{}{"this", "that"}
	ev_root.Arguments["value"] = "the other thing"
	ev_child.Arguments["path"] = []interface{}{"simple"}
	ev_child.Arguments["value"] = "and would work"

	ev_child.SetParent(ev_root)
	ev_child.Register()
	ev_root.Register()

	err := ev_child.Goto()
	if err == nil {
		t.Fatal("Goto with unapplyable parent should fail!")
	}
}
