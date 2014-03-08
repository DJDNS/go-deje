package state

import (
	"reflect"
	"testing"
)

func TestGetTraversal(t *testing.T) {
	original := map[string]interface{}{
		"hello": "world",
		"deep": map[string]interface{}{
			"stuff": "in here",
		},
	}
	container, err := MakeContainer(original)
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = getTraversal(container, []interface{}{})
	if err == nil {
		t.Fatal("getTraversal should fail for zero-length path")
	}

	_, _, err = getTraversal(container, []interface{}{"this", "that"})
	if err == nil {
		t.Fatal("getTraversal should fail for bad traversal")
	}

	p, l, err := getTraversal(container, []interface{}{"this"})
	if err != nil {
		t.Error("getTraversal should not fail for nonexistent last")
		t.Fatal(err)
	}
	expected := original
	if !reflect.DeepEqual(p.Export(), expected) {
		t.Fatalf("Expected %#v, got %#v", expected, p.Export())
	}
	if l != "this" {
		t.Fatalf("Expected %#v, got %#v", "this", l)
	}

	p, l, err = getTraversal(container, []interface{}{"deep", "stuff"})
	if err != nil {
		t.Error("getTraversal should not fail for deep traversal")
		t.Fatal(err)
	}
	expected = map[string]interface{}{"stuff": "in here"}
	if !reflect.DeepEqual(p.Export(), expected) {
		t.Fatalf("Expected %#v, got %#v", expected, p.Export())
	}
	if l != "stuff" {
		t.Fatalf("Expected %#v, got %#v", "this", l)
	}
}

type primitiveApplyTest struct {
	Change     Primitive
	Expected   interface{}
	FailureMsg string
}

func (test *primitiveApplyTest) Run(t *testing.T, ds *DocumentState) {
	err := test.Change.Apply(ds)
	if err != nil {
		t.Fatal(err)
	}
	exported := ds.Export()
	if !reflect.DeepEqual(exported, test.Expected) {
		t.Error(test.FailureMsg)
		t.Fatalf("Expected %#v, got %#v", test.Expected, exported)
	}
}

func TestDeletePrimitive_Apply_Root(t *testing.T) {
	ds := NewDocumentState()
	primitive := DeletePrimitive{
		Path: []interface{}{},
	}
	err := primitive.Apply(ds)
	if err == nil {
		t.Fatal("DeletePrimitive.Apply should always fail for root")
	}
}

func TestDeletePrimitive_Apply_WithPath(t *testing.T) {
	ds := NewDocumentState()
	primitive := DeletePrimitive{
		Path: []interface{}{"this", "that"},
	}
	err := primitive.Apply(ds)
	if err == nil {
		t.Fatal("DeletePrimitive.Apply should fail for bad traversals")
	}
	primitive.Path = []interface{}{"this"}
	err = primitive.Apply(ds)
	if err != nil {
		t.Fatal(err) // Should not fail - deletion is idempotent
	}

	// Confirm that actual deletion happens
	setter := SetPrimitive{
		Path: []interface{}{},
		Value: map[string]interface{}{
			"hello": "world",
			"other": "stuff",
			"deep": map[string]interface{}{
				"this":   "that",
				"things": []interface{}{true, false, nil},
			},
		},
	}
	err = setter.Apply(ds)
	if err != nil {
		t.Fatal(err)
	}

	tests := []primitiveApplyTest{
		primitiveApplyTest{
			&DeletePrimitive{[]interface{}{"this"}},
			map[string]interface{}{
				"hello": "world",
				"other": "stuff",
				"deep": map[string]interface{}{
					"this":   "that",
					"things": []interface{}{true, false, nil},
				},
			},
			"Deleting nonexistent key should have no effect",
		},
		primitiveApplyTest{
			&DeletePrimitive{[]interface{}{"other"}},
			map[string]interface{}{
				"hello": "world",
				"deep": map[string]interface{}{
					"this":   "that",
					"things": []interface{}{true, false, nil},
				},
			},
			"Should have removed 'other' item from DS",
		},
		primitiveApplyTest{
			&DeletePrimitive{[]interface{}{"deep", "things", 1}},
			map[string]interface{}{
				"hello": "world",
				"deep": map[string]interface{}{
					"this":   "that",
					"things": []interface{}{true, nil},
				},
			},
			"Should have removed deep-traversal item from DS",
		},
	}
	for _, test := range tests {
		test.Run(t, ds)
	}
}

type primitiveReverseTest struct {
	Original   interface{}
	Change     Primitive
	FailureMsg string
}

func (test *primitiveReverseTest) Run(t *testing.T) {
	ds := NewDocumentState()
	setter := SetPrimitive{
		Path:  []interface{}{},
		Value: test.Original,
	}
	err := setter.Apply(ds)
	if err != nil {
		t.Error(test.FailureMsg)
		t.Fatal(err)
	}

	reverse, err := test.Change.Reverse(ds)
	if err != nil {
		t.Error(test.FailureMsg)
		t.Fatal(err)
	}

	err = test.Change.Apply(ds)
	if err != nil {
		t.Error(test.FailureMsg)
		t.Fatal(err)
	}

	err = reverse.Apply(ds)
	if err != nil {
		t.Error(test.FailureMsg)
		t.Fatal(err)
	}
	exported := ds.Export()
	if !reflect.DeepEqual(exported, test.Original) {
		t.Error(test.FailureMsg)
		t.Fatalf("Expected %#v, got %#v", test.Original, exported)
	}
}

func TestDeletePrimitive_Reverse(t *testing.T) {
	tests := []primitiveReverseTest{
		primitiveReverseTest{
			Original: map[string]interface{}{},
			Change: &DeletePrimitive{
				Path: []interface{}{"missing"},
			},
			FailureMsg: "Reversal on root change",
		},
		primitiveReverseTest{
			Original: map[string]interface{}{
				"existing": "stuff",
			},
			Change: &DeletePrimitive{
				Path: []interface{}{"existing"},
			},
			FailureMsg: "Reversal on root restores old root",
		},
		primitiveReverseTest{
			Original: map[string]interface{}{
				"existing": "stuff",
				"hello":    "world",
			},
			Change: &DeletePrimitive{
				Path: []interface{}{"hello"},
			},
			FailureMsg: "Reversal on root does not affect existing",
		},
	}

	for _, test := range tests {
		test.Run(t)
	}
}
