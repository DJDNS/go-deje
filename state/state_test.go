package state

import (
	"reflect"
	"testing"
)

func TestNewDocumentState(t *testing.T) {
	ds := NewDocumentState()
	if ds.Value == nil {
		t.Fatal("ds.Value is not supposed to be nil")
	}

	expected := make(map[string]interface{})
	exported := ds.Value.Export()
	if !reflect.DeepEqual(exported, expected) {
		t.Fatal("Expected %#v, got %#v", expected, exported)
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
