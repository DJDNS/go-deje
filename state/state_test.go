package state

import (
	"reflect"
	"testing"
	"time"
)

func TestNewDocumentState(t *testing.T) {
	ds := NewDocumentState()
	if ds.Value == nil {
		t.Fatal("ds.Value is not supposed to be nil")
	}

	expected := make(map[string]interface{})
	exported := ds.Value.Export()
	if !reflect.DeepEqual(exported, expected) {
		t.Fatalf("Expected %#v, got %#v", expected, exported)
	}
}

func TestDocumentState_Reset(t *testing.T) {
	ds := NewDocumentState()
	setter := &SetPrimitive{
		Path: []interface{}{},
		Value: map[string]interface{}{
			"hello": "world",
		},
	}
	err := setter.Apply(ds)
	if err != nil {
		t.Fatal(err)
	}

	expected := setter.Value
	exported := ds.Export()
	if !reflect.DeepEqual(exported, expected) {
		t.Fatalf("Expected %#v, got %#v", expected, exported)
	}

	sub := ds.Subscribe()
	ds.Reset()

	expected = map[string]interface{}{}
	exported = ds.Export()
	if !reflect.DeepEqual(exported, expected) {
		t.Fatalf("Expected %#v, got %#v", expected, exported)
	}

	select {
	case primitive := <-sub.Out():
		expected_p := &SetPrimitive{
			Path:  []interface{}{},
			Value: map[string]interface{}{},
		}
		if !reflect.DeepEqual(primitive, expected_p) {
			t.Fatalf("Expected %#v, got %#v", expected_p, primitive)
		}
		if sub.Len() != 0 {
			t.Errorf("sub should be empty, still %d left", sub.Len())
		}
	case <-time.After(time.Millisecond):
		t.Fatal("No primitive received")
	}
}

func TestDocumentState_Apply(t *testing.T) {
	ds := NewDocumentState()
	primitive := &SetPrimitive{
		Path:  []interface{}{"key"},
		Value: "value",
	}
	sub := ds.Subscribe()
	err := ds.Apply(primitive)
	if err != nil {
		t.Fatal(err)
	}
	if sub.Len() != 1 {
		t.Fatal(
			"Expected 1 primitive to be broadast, got %d",
			sub.Len(),
		)
	}
	recvd_p := <-sub.Out()
	if !reflect.DeepEqual(recvd_p, primitive) {
		t.Fatal("Expected %#v, got %#v", primitive, recvd_p)
	}
}
func TestDocumentState_Apply_BadPrimitive(t *testing.T) {
	ds := NewDocumentState()
	primitive := &SetPrimitive{
		Path:  []interface{}{"no", "such", "path"},
		Value: 8,
	}
	sub := ds.Subscribe()
	err := ds.Apply(primitive)
	if err == nil {
		t.Fatal("ds.Apply should fail if underlying Apply fails")
	}
	if sub.Len() != 0 {
		t.Fatal(
			"Expected 0 primitives to be broadast, got %d",
			sub.Len(),
		)
	}
}

func TestDocumentState_Export(t *testing.T) {
	ds := NewDocumentState()
	err := ds.Value.SetChild("hello", "world")
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]interface{}{
		"hello": "world",
	}
	exported := ds.Export()
	if !reflect.DeepEqual(exported, expected) {
		t.Fatal("Expected %#v, got %#v", expected, exported)
	}
}
