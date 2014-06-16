package deje

import (
	"github.com/campadrenalin/go-deje/state"
	"reflect"
	"testing"
)

func TestSimpleClient_NewSimpleClient(t *testing.T) {
	topic := "http://example.com/deje/some-doc"
	sc := NewSimpleClient(topic)
	if sc.client.Doc.Topic != topic {
		t.Fatal("Did not create encapsulated Client correctly")
	}
}

func TestSimpleClient_Connect(t *testing.T) {
	topic := "http://example.com/deje/some-doc"
	client := NewSimpleClient(topic)
	server_addr, server_closer := setupServer()
	defer server_closer()

	err := client.Connect("foo")
	if err == nil {
		t.Fatal("foo is not a real server - should not 'succeed'")
	}

	err = client.Connect(server_addr)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: Test that RequestTip was broadcast
}

func TestSimpleClient_GetDoc(t *testing.T) {
	topic := "http://example.com/deje/some-doc"
	client := NewSimpleClient(topic)

	got := client.GetDoc()
	expected := client.client.Doc
	if got != expected {
		t.Fatal("GetDoc returned wrong pointer")
	}
}

func TestSimpleClient_Export(t *testing.T) {
	topic := "http://example.com/deje/some-doc"
	client := NewSimpleClient(topic)

	// Test before any changes
	exported := client.Export()
	expected := map[string]interface{}{}
	if !reflect.DeepEqual(exported, expected) {
		t.Fatalf("Expected %#v, got %#v", expected, exported)
	}

	// Update the contents of the Doc
	primitive := state.SetPrimitive{
		Path: []interface{}{},
		Value: map[string]interface{}{
			"Rabbit": "rabbit",
		},
	}
	client.client.Doc.State.Apply(&primitive)

	// Test that the new contents reflect the changes
	exported = client.Export()
	expected["Rabbit"] = "rabbit"
	if !reflect.DeepEqual(exported, expected) {
		t.Fatalf("Expected %#v, got %#v", expected, exported)
	}
}
