package deje

import (
	"bytes"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/campadrenalin/go-deje/state"
	"github.com/stretchr/testify/assert"
)

func TestSimpleClient_NewSimpleClient(t *testing.T) {
	topic := "http://example.com/deje/some-doc"
	sc := NewSimpleClient(topic, nil)
	if sc.client.Doc.Topic != topic {
		t.Fatal("Did not create encapsulated Client correctly")
	}
}

func TestSimpleClient_Connect(t *testing.T) {
	topic := "http://example.com/deje/some-doc"
	client := NewSimpleClient(topic, nil)
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
	Simple     []*SimpleClient
	Logs       []*bytes.Buffer
	Listener   Client
	EventsRcvd chan interface{}
	Closer     func()
}

func setupSimpleProtocolTest(t *testing.T, num_simple int) simpleProtoTest {
	var spt simpleProtoTest
	spt.Topic = "http://example.com/deje/some-doc"
	spt.Simple = make([]*SimpleClient, num_simple)
	spt.Logs = make([]*bytes.Buffer, num_simple)
	spt.Listener = NewClient(spt.Topic)
	server_addr, server_closer := setupServer()
	spt.Closer = server_closer

	// Use this order to ignore any RequestTip() called during Connect()
	spt.EventsRcvd = make(chan interface{}, 10)
	for i := 0; i < num_simple; i++ {
		buffer := new(bytes.Buffer)
		logger := log.New(buffer, "deje.SimpleClient: ", 0)

		spt.Logs[i] = buffer
		spt.Simple[i] = NewSimpleClient(spt.Topic, logger)
		if err := spt.Simple[i].Connect(server_addr); err != nil {
			t.Fatal(err)
		}
	}
	if err := spt.Listener.Connect(server_addr); err != nil {
		t.Fatal(err)
	}

	// Make sure all connect fully, THEN start listening
	<-time.After(50 * time.Millisecond)
	spt.Listener.SetEventCallback(func(event interface{}) {
		spt.EventsRcvd <- event
	})

	return spt
}

func (spt simpleProtoTest) Expect(t *testing.T, messages []interface{}) {
	for id, expected := range messages {
		select {
		case event := <-spt.EventsRcvd:
			if !reflect.DeepEqual(event, expected) {
				t.Fatalf("\nexp %#v\ngot %#v", expected, event)
			}
		case <-time.After(50 * time.Millisecond):
			t.Fatalf("Timed out waiting for event %d (%#v)", id, expected)
		}
	}
	// Ensure no extra events after
	<-time.After(5 * time.Millisecond)
	if len(spt.EventsRcvd) != 0 {
		t.Fatal("Wrong number of events received")
	}
}

func TestSimpleClient_RequestTip(t *testing.T) {
	spt := setupSimpleProtocolTest(t, 1)
	defer spt.Closer()

	if err := spt.Simple[0].RequestTip(); err != nil {
		t.Fatal(err)
	}
	spt.Expect(t, []interface{}{
		map[string]interface{}{
			"type": "01-request-tip",
		},
	})
}

func TestSimpleClient_PublishTip(t *testing.T) {
	spt := setupSimpleProtocolTest(t, 1)
	defer spt.Closer()

	spt.Simple[0].tip = "some hash"
	if err := spt.Simple[0].PublishTip(); err != nil {
		t.Fatal(err)
	}
	spt.Expect(t, []interface{}{
		map[string]interface{}{
			"type":     "01-publish-tip",
			"tip_hash": "some hash",
		},
	})
}

func TestSimpleClient_TipCycle(t *testing.T) {
	spt := setupSimpleProtocolTest(t, 2)
	defer spt.Closer()

	spt.Simple[0].tip = "some hash" // Make sure requesting client does not ask for history
	spt.Simple[1].tip = "some hash"
	if err := spt.Simple[0].RequestTip(); err != nil {
		t.Fatal(err)
	}
	spt.Expect(t, []interface{}{
		map[string]interface{}{
			"type": "01-request-tip",
		},
		map[string]interface{}{
			"type":     "01-publish-tip",
			"tip_hash": "some hash",
		},
	})
}

func TestSimpleClient_Rcv_BadMsg(t *testing.T) {
	spt := setupSimpleProtocolTest(t, 2)
	defer spt.Closer()

	// Send a series of bad data
	// (can't do numbers, floating point eq fails)
	messages := []interface{}{
		"Not a map, muahaha", true, false, nil,
		[]interface{}{}, []interface{}{"x", "y", "z"},
		map[string]interface{}{
			"type": true,
		},
		map[string]interface{}{
			"type": "foo",
		},
		map[string]interface{}{
			"no_type_key": "frowny face",
		},
		map[string]interface{}{},
	}
	for _, msg := range messages {
		if err := spt.Simple[0].client.Publish(msg); err != nil {
			t.Fatal(err)
		}
	}
	_unf_msg_type := "deje.SimpleClient: Unfamiliar message type: "
	_non_obj_msg := "deje.SimpleClient: Non-{} message\n"
	_no_type_parm := "deje.SimpleClient: Message with no 'type' param\n"
	expected_log := _non_obj_msg +
		_non_obj_msg +
		_non_obj_msg +
		_non_obj_msg +
		_non_obj_msg +
		_non_obj_msg +
		_no_type_parm +
		_unf_msg_type + "'foo'\n" +
		_no_type_parm +
		_no_type_parm

	// Expect only evil data, no error
	spt.Expect(t, messages)
	assert.Equal(t, expected_log, spt.Logs[1].String())

	// Confirm that we still respond well to legit data afterwards
	spt.Simple[0].tip = "some hash" // Make sure requesting client does not ask for history
	spt.Simple[1].tip = "some hash"
	if err := spt.Simple[0].RequestTip(); err != nil {
		t.Fatal(err)
	}
	spt.Expect(t, []interface{}{
		map[string]interface{}{
			"type": "01-request-tip",
		},
		map[string]interface{}{
			"type":     "01-publish-tip",
			"tip_hash": "some hash",
		},
	})
}

func TestSimpleClient_RequestHistory(t *testing.T) {
	spt := setupSimpleProtocolTest(t, 1)
	defer spt.Closer()

	if err := spt.Simple[0].RequestHistory(); err != nil {
		t.Fatal(err)
	}
	spt.Expect(t, []interface{}{
		map[string]interface{}{
			"type": "01-request-history",
		},
	})
}

func TestSimpleClient_PublishHistory_NoHistory(t *testing.T) {
	spt := setupSimpleProtocolTest(t, 1)
	defer spt.Closer()

	if err := spt.Simple[0].PublishHistory(); err != nil {
		t.Fatal(err)
	}
	spt.Expect(t, []interface{}{
		map[string]interface{}{
			"type":     "01-publish-history",
			"tip_hash": "",
			"error":    "not-found",
		},
	})
}

func TestSimpleClient_PublishHistory_IncompleteHistory(t *testing.T) {
	spt := setupSimpleProtocolTest(t, 1)
	defer spt.Closer()

	// Add a bit of history, but never register root
	doc := spt.Simple[0].GetDoc()
	root := doc.NewEvent("handler name")
	child := doc.NewEvent("other handler")
	child.SetParent(root)
	child.Register()
	spt.Simple[0].tip = child.Hash()

	if err := spt.Simple[0].PublishHistory(); err != nil {
		t.Fatal(err)
	}
	spt.Expect(t, []interface{}{
		map[string]interface{}{
			"type":     "01-publish-history",
			"tip_hash": child.Hash(),
			"error":    "root-not-found",
		},
	})
}

func TestSimpleClient_PublishHistory_FullHistory(t *testing.T) {
	spt := setupSimpleProtocolTest(t, 1)
	defer spt.Closer()

	// Add full history
	doc := spt.Simple[0].GetDoc()
	root := doc.NewEvent("root")
	child := doc.NewEvent("child")
	child.SetParent(root)
	root.Register()
	child.Register()
	spt.Simple[0].tip = child.Hash()

	if err := spt.Simple[0].PublishHistory(); err != nil {
		t.Fatal(err)
	}
	spt.Expect(t, []interface{}{
		map[string]interface{}{
			"type":     "01-publish-history",
			"tip_hash": child.Hash(),
			"history": []interface{}{
				map[string]interface{}{
					"handler": "root",
					"parent":  "",
					"args":    map[string]interface{}{},
				},
				map[string]interface{}{
					"handler": "child",
					"parent":  root.Hash(),
					"args":    map[string]interface{}{},
				},
			},
		},
	})
}

// Can test a simple failure, because we cover (more complex) success
// in other tests already.
func TestSimpleClient_HistoryCycle(t *testing.T) {
	spt := setupSimpleProtocolTest(t, 2)
	defer spt.Closer()

	spt.Simple[1].tip = "some hash"
	if err := spt.Simple[0].RequestHistory(); err != nil {
		t.Fatal(err)
	}
	spt.Expect(t, []interface{}{
		map[string]interface{}{
			"type": "01-request-history",
		},
		map[string]interface{}{
			"type":     "01-publish-history",
			"tip_hash": "some hash",
			"error":    "not-found",
		},
	})
}

func TestSimpleClient_Promote(t *testing.T) {
	spt := setupSimpleProtocolTest(t, 2)
	defer spt.Closer()

	doc1 := spt.Simple[0].GetDoc()
	doc2 := spt.Simple[1].GetDoc()

	event := doc1.NewEvent("SET")
	if err := spt.Simple[0].Promote(event); err == nil {
		t.Fatal("Should fail if we can't navigate to event!")
	}

	event.Arguments["path"] = []interface{}{"bar"}
	event.Arguments["value"] = "baz"
	event.Register()

	if err := spt.Simple[0].Promote(event); err != nil {
		t.Fatal(err)
	}
	spt.Expect(t, []interface{}{
		map[string]interface{}{
			"type":     "01-publish-tip",
			"tip_hash": event.Hash(),
		},
		map[string]interface{}{
			"type": "01-request-history",
		},
		map[string]interface{}{
			"type":     "01-publish-history",
			"tip_hash": event.Hash(),
			"history": []interface{}{
				map[string]interface{}{
					"handler": "SET",
					"parent":  "",
					"args":    event.Arguments,
				},
			},
		},
	})

	assert.Equal(t, spt.Simple[0].tip, event.Hash())
	assert.Equal(t, spt.Simple[1].tip, event.Hash())
	assert.Equal(t, *doc1.Events[event.Hash()], event)
	assert.Equal(t, *doc2.Events[event.Hash()], event)

	expected_export := map[string]interface{}{
		"bar": "baz",
	}
	assert.Equal(t, spt.Simple[0].Export(), expected_export)
	assert.Equal(t, spt.Simple[1].Export(), expected_export)
}

func TestSimpleClient_GetDoc(t *testing.T) {
	topic := "http://example.com/deje/some-doc"
	client := NewSimpleClient(topic, nil)

	got := client.GetDoc()
	expected := client.client.Doc
	if got != expected {
		t.Fatal("GetDoc returned wrong pointer")
	}
}

func TestSimpleClient_Export(t *testing.T) {
	topic := "http://example.com/deje/some-doc"
	client := NewSimpleClient(topic, nil)

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
