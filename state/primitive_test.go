package state

import (
	"reflect"
	"testing"
)

func TestSetPrimitive_Apply_Root(t *testing.T) {
	ds := NewDocumentState()
	primitive := SetPrimitive{
		Path:  []interface{}{},
		Value: map[string]interface{}{"hello": "world"},
	}
	err := primitive.Apply(ds)
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]interface{}{
		"hello": "world",
	}
	exported := ds.Export()
	if !reflect.DeepEqual(exported, expected) {
		t.Fatalf("Expected %#v, got %#v", expected, exported)
	}

	primitive.Value.(map[string]interface{})["invalid"] = make(chan int)
	err = primitive.Apply(ds)
	if err == nil {
		t.Fatal("primitive.Apply should fail when given bad value")
	}
}

func TestSetPrimitive_Apply_WithPath(t *testing.T) {
	// Set up initial state
	ds := NewDocumentState()
	primitive := SetPrimitive{
		Path: []interface{}{},
		Value: map[string]interface{}{
			"deep": []interface{}{
				"stuff", "in", "here",
			},
		},
	}
	err := primitive.Apply(ds)
	if err != nil {
		t.Fatal(err)
	}

	// Try to reapply it deeper, Inception-style.
	// But first, let's be stupid, make sure things don't shatter.
	primitive.Path = []interface{}{true}
	if primitive.Apply(ds) == nil {
		t.Fatal("primitive.Apply should fail for bad key type")
	}
	primitive.Path = []interface{}{"non-existent key", 0}
	if primitive.Apply(ds) == nil {
		t.Fatal("primitive.Apply should fail for bad traversal")
	}
	primitive.Path = []interface{}{"deep", "wide"}
	if primitive.Apply(ds) == nil {
		t.Fatal("primitive.Apply should fail for bad traversal")
	}

	// Alright, let's get on with it.
	primitive.Path = []interface{}{"deep", 2}
	err = primitive.Apply(ds)
	if err != nil {
		t.Fatal(err)
	}

	// Confirm it works
	expected := map[string]interface{}{
		"deep": []interface{}{
			"stuff", "in",
			map[string]interface{}{
				"deep": []interface{}{
					"stuff", "in", "here",
				},
			},
		},
	}
	exported := ds.Export()
	if !reflect.DeepEqual(exported, expected) {
		t.Fatalf("Expected %#v, got %#v", expected, exported)
	}
}
