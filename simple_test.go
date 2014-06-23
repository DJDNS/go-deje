package deje

import (
	"github.com/campadrenalin/go-deje/state"
	"reflect"
	"testing"
	"time"
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
	listener := NewClient(topic)
	server_addr, server_closer := setupServer()
	defer server_closer()

	err := client.Connect("foo")
	if err == nil {
		t.Fatal("foo is not a real server - should not 'succeed'")
	}

	// Set up listener to detect initial RequestTip
	events_rcvd := make(chan interface{}, 10)
	listener.SetEventCallback(func(event interface{}) {
		events_rcvd <- event
	})
	if err := listener.Connect(server_addr); err != nil {
		t.Fatal(err)
	}

	// Connect the SimpleClient
	err = client.Connect(server_addr)
	if err != nil {
		t.Fatal(err)
	}

	// Ensure that RequestTip was broadcast
	expected := map[string]interface{}{
		"type": "01-request-tip",
	}
	select {
	case event := <-events_rcvd:
		if !reflect.DeepEqual(event, expected) {
			t.Fatalf("Expected %#v, got %#v", expected, event)
		}
	case <-time.After(50 * time.Millisecond):
		t.Fatal("Timed out waiting for event")
	}
	// Ensure no extra events after
	if len(events_rcvd) != 0 {
		t.Fatal("Wrong number of events received")
	}
}

type simpleProtoTest struct {
	Topic      string
	Simple     SimpleClient
	Listener   Client
	EventsRcvd chan interface{}
	Closer     func()
}

func setupSimpleProtocolTest(t *testing.T) simpleProtoTest {
	var spt simpleProtoTest
	spt.Topic = "http://example.com/deje/some-doc"
	spt.Simple = NewSimpleClient(spt.Topic)
	spt.Listener = NewClient(spt.Topic)
	server_addr, server_closer := setupServer()
	spt.Closer = server_closer

	// Use this order to ignore any RequestTip() called during Connect()
	spt.EventsRcvd = make(chan interface{}, 10)
	if err := spt.Simple.Connect(server_addr); err != nil {
		t.Fatal(err)
	}
	if err := spt.Listener.Connect(server_addr); err != nil {
		t.Fatal(err)
	}
	spt.Listener.SetEventCallback(func(event interface{}) {
		spt.EventsRcvd <- event
	})
	<-time.After(50 * time.Millisecond) // Make sure both connect fully

	return spt
}

func TestSimpleClient_RequestTip(t *testing.T) {
	spt := setupSimpleProtocolTest(t)
	defer spt.Closer()

	expected := map[string]interface{}{
		"type": "01-request-tip",
	}
	if err := spt.Simple.RequestTip(); err != nil {
		t.Fatal(err)
	}
	select {
	case event := <-spt.EventsRcvd:
		if !reflect.DeepEqual(event, expected) {
			t.Fatalf("Expected %#v, got %#v", expected, event)
		}
	case <-time.After(50 * time.Millisecond):
		t.Fatal("Timed out waiting for event")
	}
	// Ensure no extra events after
	if len(spt.EventsRcvd) != 0 {
		t.Fatal("Wrong number of events received")
	}
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
