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
