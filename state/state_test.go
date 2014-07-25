package state

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

	primitives_applied := make(chan Primitive, 10)
	ds.SetPrimitiveCallback(func(p Primitive) {
		primitives_applied <- p
	})
	ds.Reset()

	expected = map[string]interface{}{}
	exported = ds.Export()
	if !reflect.DeepEqual(exported, expected) {
		t.Fatalf("Expected %#v, got %#v", expected, exported)
	}

	select {
	case primitive := <-primitives_applied:
		expected_p := &SetPrimitive{
			Path:  []interface{}{},
			Value: map[string]interface{}{},
		}
		if !reflect.DeepEqual(primitive, expected_p) {
			t.Fatalf("Expected %#v, got %#v", expected_p, primitive)
		}
		num_left := len(primitives_applied)
		if num_left != 0 {
			t.Errorf(
				"primitives_applied should be empty, still %d left",
				num_left,
			)
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
	primitives_applied := make(chan Primitive, 10)
	ds.SetPrimitiveCallback(func(p Primitive) {
		primitives_applied <- p
	})
	err := ds.Apply(primitive)
	if err != nil {
		t.Fatal(err)
	}
	if len(primitives_applied) != 1 {
		t.Fatal(
			"Expected 1 primitive to be broadast, got %d",
			len(primitives_applied),
		)
	}
	recvd_p := <-primitives_applied
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
	primitives_applied := make(chan Primitive, 10)
	ds.SetPrimitiveCallback(func(p Primitive) {
		primitives_applied <- p
	})
	err := ds.Apply(primitive)
	if err == nil {
		t.Fatal("ds.Apply should fail if underlying Apply fails")
	}
	if len(primitives_applied) != 0 {
		t.Fatal(
			"Expected 0 primitives to be broadast, got %d",
			len(primitives_applied),
		)
	}
}
func TestDocumentState_Apply_NilCallback(t *testing.T) {
	ds := NewDocumentState()
	primitive := &SetPrimitive{
		Path:  []interface{}{"key"},
		Value: "value",
	}
	ds.SetPrimitiveCallback(nil)

	// Should just work with no panics or errors
	if err := ds.Apply(primitive); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, map[string]interface{}{"key": "value"}, ds.Export())
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
