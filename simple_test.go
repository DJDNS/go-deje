package deje

import (
	"bytes"
	"errors"
	"log"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/DJDNS/go-deje/state"
	"github.com/stretchr/testify/assert"
)

func TestSimpleClient_NewSimpleClient(t *testing.T) {
	topic := "http://example.com/deje/some-doc"
	sc := NewSimpleClient(topic, nil)
	if sc.client.Topic != topic {
		t.Fatal("Did not create encapsulated Client correctly")
	}
}

func TestSimpleClient_Open_BadUrl(t *testing.T) {
	_, err := Open("localhost:8080", nil, nil)
	if assert.Error(t, err, "Open should have failed, due to bad URL") {
		assert.Equal(t, err.Error(), "URL does not start with 'deje://': 'localhost:8080'")
	}
}

func TestSimpleClient_Open_NoSuchHost(t *testing.T) {
	_, err := Open("deje://no.such.host:8080/", nil, nil)
	if assert.Error(t, err, "Open should have failed, due to unreachable host") {
		assert.Equal(t, err.Error(), "Error connecting to websocket server: websocket.Dial ws://no.such.host:8080/ws: dial tcp: lookup no.such.host: no such host")
	}
}

func TestSimpleClient_Open(t *testing.T) {
	buffer := new(bytes.Buffer)
	logger := log.New(buffer, "deje.SimpleClient: ", 0)
	server_addr, server_closer := setupServer()
	defer server_closer()

	var got_a_primitive bool
	handler := func(primitive state.Primitive) {
		_, ok := primitive.(*state.SetPrimitive)
		assert.True(t, ok, "Got a SetPrimitive")
		got_a_primitive = true
	}

	url := strings.Replace(server_addr, "ws://", "deje://", 1) + "/some/topic"
	sc, err := Open(url, logger, handler)
	if !assert.NoError(t, err, "Open should succeed for URL '%s'", url) {
		t.Fail()
	}

	topic := url // The two are equal, unless deje_url is missing path component
	if !assert.Equal(t, topic, sc.client.Topic, "Should subscribe to correct topic") {
		t.Fail()
	}

	// Set up event in SimpleClient
	event := sc.GetDoc().NewEvent("SET")
	event.Arguments["path"] = []interface{}{"message"}
	event.Arguments["value"] = "Karma incremented"
	event.Register() // But do not Goto() yet

	// raw_client used to trigger Goto over the network
	raw_client := NewClient(topic)
	if err := raw_client.Connect(server_addr); err != nil {
		t.Fatal(err)
	}
	message := map[string]interface{}{
		"type":       "02-publish-timestamps",
		"timestamps": []interface{}{event.Hash()},
	}
	if !assert.NoError(t, raw_client.Publish(message)) {
		t.Fail()
	}

	// Confirm that content has been updated
	<-time.After(50 * time.Millisecond)
	assert.Equal(t,
		map[string]interface{}{
			"message": "Karma incremented",
		},
		sc.Export(),
	)

	// Confirm that OnPrimitiveCallback was called
	assert.True(t, got_a_primitive, "Recvd a primitive, callback was called")
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
		"type": "02-request-timestamps",
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

type failingTimestampService string

func (f failingTimestampService) GetTimestamps() ([]string, error) {
	return nil, errors.New(string(f))
}

func TestSimpleClient_ReTip_Fail(t *testing.T) {
	buffer := new(bytes.Buffer)
	logger := log.New(buffer, "retip_test: ", 0)
	sc := NewSimpleClient("deje://demo/", logger)
	sc.tt.Service = failingTimestampService("Bad service for retip")

	sc.ReTip()
	assert.Equal(t, "retip_test: Bad service for retip\n", buffer.String())
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

type logtest struct {
	Message interface{}
	Logline string
}

func (lt logtest) Run(t *testing.T, spt simpleProtoTest) {
	spt.Logs[1].Reset()

	if err := spt.Simple[0].client.Publish(lt.Message); err != nil {
		t.Fatal(err)
	}
	spt.Expect(t, []interface{}{lt.Message})

	var expected_log string
	if lt.Logline != "" {
		expected_log = "deje.SimpleClient: " + lt.Logline + "\n"
	}
	assert.Equal(t, expected_log, spt.Logs[1].String())
}

func TestSimpleClient_Rcv_BadMsg(t *testing.T) {
	spt := setupSimpleProtocolTest(t, 2)
	defer spt.Closer()

	// Set up error messages for reuse
	_unf_msg_type := "Unfamiliar message type: "
	_non_obj_msg := "Non-{} message"
	_no_type_param := "Message with no 'type' param"
	_bad_events := "Message with bad 'events' param"
	_bad_ts := "Message with bad 'timestamps' param"
	_clone_err := "json: cannot unmarshal bool into Go value of type document.Event"

	// Cannot be Goto'd
	//incomplete_event := spt.Simple[0].GetDoc().NewEvent("SET")

	// Send a series of bad data
	// (can't do numbers, floating point eq fails)
	logtests := []logtest{
		logtest{
			"Not a map, muahaha",
			_non_obj_msg,
		},
		logtest{
			true,
			_non_obj_msg,
		},
		logtest{
			false,
			_non_obj_msg,
		},
		logtest{
			nil,
			_non_obj_msg,
		},
		logtest{
			[]interface{}{},
			_non_obj_msg,
		},
		logtest{
			[]interface{}{"x", "y", "z"},
			_non_obj_msg,
		},
		logtest{
			map[string]interface{}{
				"type": true,
			},
			_no_type_param,
		},
		logtest{
			map[string]interface{}{
				"type": "foo",
			},
			_unf_msg_type + "'foo'",
		},
		logtest{
			map[string]interface{}{
				"no_type_key": "frowny face",
			},
			_no_type_param,
		},
		logtest{
			map[string]interface{}{},
			_no_type_param,
		},
		logtest{
			map[string]interface{}{
				"type": "02-publish-events",
			},
			_bad_events,
		},
		logtest{
			map[string]interface{}{
				"type":   "02-publish-events",
				"events": []interface{}{true},
			},
			_clone_err,
		},
		logtest{
			map[string]interface{}{
				"type": "02-publish-timestamps",
			},
			_bad_ts,
		},
		logtest{
			map[string]interface{}{
				"type":       "02-publish-timestamps",
				"timestamps": "",
			},
			_bad_ts,
		},
		logtest{
			map[string]interface{}{
				"type":       "02-publish-timestamps",
				"timestamps": []interface{}{"Only strings allowed!", true, false},
			},
			_bad_ts,
		},
		logtest{
			map[string]interface{}{
				"type": "log",
			},
			"",
		},
	}
	for _, lt := range logtests {
		lt.Run(t, spt)
	}

	// Confirm that we still respond well to legit data afterwards
	if err := spt.Simple[0].RequestTimestamps(); err != nil {
		t.Fatal(err)
	}
	spt.Expect(t, []interface{}{
		map[string]interface{}{
			"type": "02-request-timestamps",
		},
		map[string]interface{}{
			"type":       "02-publish-timestamps",
			"timestamps": []interface{}{},
		},
	})
}

func TestSimpleClient_RequestEvents(t *testing.T) {
	spt := setupSimpleProtocolTest(t, 1)
	defer spt.Closer()

	if err := spt.Simple[0].RequestEvents(); err != nil {
		t.Fatal(err)
	}
	spt.Expect(t, []interface{}{
		map[string]interface{}{
			"type": "02-request-events",
		},
	})
}

func TestSimpleClient_PublishEvents(t *testing.T) {
	spt := setupSimpleProtocolTest(t, 1)
	defer spt.Closer()

	if err := spt.Simple[0].PublishEvents(); err != nil {
		t.Fatal(err)
	}
	spt.Expect(t, []interface{}{
		map[string]interface{}{
			"type":   "02-publish-events",
			"events": []interface{}{},
		},
	})
}

func TestSimpleClient_EventCycle(t *testing.T) {
	spt := setupSimpleProtocolTest(t, 2)
	defer spt.Closer()

	doc1 := spt.Simple[0].GetDoc()
	doc2 := spt.Simple[1].GetDoc()

	first_event := doc1.NewEvent("first")
	first_event.Arguments["nonce"] = "00"
	first_event.Register()
	second_event := doc1.NewEvent("second")
	second_event.Register()

	assert.True(t,
		first_event.Hash() < second_event.Hash(),
		"Events sort according to their names",
		first_event.Hash(),
		second_event.Hash(),
	)

	if err := spt.Simple[1].RequestEvents(); err != nil {
		t.Fatal(err)
	}
	spt.Expect(t, []interface{}{
		map[string]interface{}{
			"type": "02-request-events",
		},
		map[string]interface{}{
			"type": "02-publish-events",
			"events": []interface{}{
				map[string]interface{}{
					"handler": "first",
					"parent":  "",
					"args": map[string]interface{}{
						"nonce": "00",
					},
				},
				map[string]interface{}{
					"handler": "second",
					"parent":  "",
					"args":    map[string]interface{}{},
				},
			},
		},
	})

	// Ensure that events were copied over
	if !assert.Equal(t, 2, len(doc2.Events)) {
		t.FailNow()
	}
	assert.Equal(t,
		doc2.Events[first_event.Hash()].HandlerName,
		first_event.HandlerName,
	)
}

func TestSimpleClient_RequestTimestamps(t *testing.T) {
	spt := setupSimpleProtocolTest(t, 1)
	defer spt.Closer()

	if err := spt.Simple[0].RequestTimestamps(); err != nil {
		t.Fatal(err)
	}
	spt.Expect(t, []interface{}{
		map[string]interface{}{
			"type": "02-request-timestamps",
		},
	})
}

func TestSimpleClient_PublishTimestamps(t *testing.T) {
	spt := setupSimpleProtocolTest(t, 1)
	defer spt.Closer()

	if err := spt.Simple[0].PublishTimestamps(); err != nil {
		t.Fatal(err)
	}
	spt.Expect(t, []interface{}{
		map[string]interface{}{
			"type":       "02-publish-timestamps",
			"timestamps": []interface{}{},
		},
	})

	doc := spt.Simple[0].GetDoc()
	doc.Timestamps = append(doc.Timestamps, "a hash", "another hash")
	if err := spt.Simple[0].PublishTimestamps(); err != nil {
		t.Fatal(err)
	}
	spt.Expect(t, []interface{}{
		map[string]interface{}{
			"type":       "02-publish-timestamps",
			"timestamps": []interface{}{"a hash", "another hash"},
		},
	})
}

func TestSimpleClient_TimestampCycle(t *testing.T) {
	spt := setupSimpleProtocolTest(t, 2)
	defer spt.Closer()

	doc0 := spt.Simple[0].GetDoc()
	doc1 := spt.Simple[1].GetDoc()

	evFirst := doc1.NewEvent("first")
	evFirst.Register()
	evSecond := doc1.NewEvent("second")
	evSecond.Register()

	expected_timestamps := []string{evFirst.Hash(), evSecond.Hash()}
	doc1.Timestamps = expected_timestamps

	// First time, events are foreign to doc0
	if err := spt.Simple[0].RequestTimestamps(); err != nil {
		t.Fatal(err)
	}
	spt.Expect(t, []interface{}{
		map[string]interface{}{
			"type": "02-request-timestamps",
		},
		map[string]interface{}{
			"type":       "02-publish-timestamps",
			"timestamps": []interface{}{evFirst.Hash(), evSecond.Hash()},
		},
		map[string]interface{}{
			"type": "02-request-events",
		},
		map[string]interface{}{
			"type": "02-publish-events",
			"events": []interface{}{
				map[string]interface{}{
					"handler": "second",
					"parent":  "",
					"args":    map[string]interface{}{},
				},
				map[string]interface{}{
					"parent":  "",
					"args":    map[string]interface{}{},
					"handler": "first",
				},
			},
		},
	})

	assert.Equal(t, expected_timestamps, doc0.Timestamps)
	assert.Equal(t, expected_timestamps, doc1.Timestamps)

	// Second time, events are known to doc0, no need to request events
	if err := spt.Simple[0].RequestTimestamps(); err != nil {
		t.Fatal(err)
	}
	spt.Expect(t, []interface{}{
		map[string]interface{}{
			"type": "02-request-timestamps",
		},
		map[string]interface{}{
			"type":       "02-publish-timestamps",
			"timestamps": []interface{}{evFirst.Hash(), evSecond.Hash()},
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
			"type":       "02-publish-timestamps",
			"timestamps": []interface{}{event.Hash()},
		},
		map[string]interface{}{
			"type": "02-request-events",
		},
		map[string]interface{}{
			"type": "02-publish-events",
			"events": []interface{}{
				map[string]interface{}{
					"handler": "SET",
					"parent":  "",
					"args":    event.Arguments,
				},
			},
		},
	})

	hash := event.Hash()
	assert.Equal(t, spt.Simple[0].Tip.Hash(), hash)
	assert.Equal(t, spt.Simple[1].Tip.Hash(), hash)
	assert.Equal(t, doc1.Events[hash].Hash(), hash)
	assert.Equal(t, doc2.Events[hash].Hash(), hash)

	expected_export := map[string]interface{}{
		"bar": "baz",
	}
	assert.Equal(t, spt.Simple[0].Export(), expected_export)
	assert.Equal(t, spt.Simple[1].Export(), expected_export)

	expected_timestamps := []string{event.Hash()}
	assert.Equal(t, doc1.Timestamps, expected_timestamps)
	assert.Equal(t, doc2.Timestamps, expected_timestamps)
}

func TestSimpleClient_SetPrimitiveCallback(t *testing.T) {
	spt := setupSimpleProtocolTest(t, 2)
	defer spt.Closer()

	primitives := make(chan state.Primitive, 10)
	on_primitive := func(p state.Primitive) {
		primitives <- p
	}
	spt.Simple[1].SetPrimitiveCallback(on_primitive)

	doc := spt.Simple[0].GetDoc()
	eventA := doc.NewEvent("SET")
	eventA.Arguments["path"] = []interface{}{"items"}
	eventA.Arguments["value"] = map[string]interface{}{
		"first":  "thing",
		"second": "thang",
	}
	eventA.Register()

	eventB := doc.NewEvent("DELETE")
	eventB.Arguments["path"] = []interface{}{"items", "second"}
	eventB.SetParent(eventA)
	eventB.Register()

	if err := spt.Simple[0].Promote(eventB); err != nil {
		t.Fatal(err)
	}
	spt.Expect(t, []interface{}{
		map[string]interface{}{
			"type":       "02-publish-timestamps",
			"timestamps": []interface{}{eventB.Hash()},
		},
		map[string]interface{}{
			"type": "02-request-events",
		},
		map[string]interface{}{
			"type": "02-publish-events",
			"events": []interface{}{
				map[string]interface{}{
					"handler": "DELETE",
					"parent":  eventA.Hash(),
					"args":    eventB.Arguments,
				},
				map[string]interface{}{
					"handler": "SET",
					"parent":  "",
					"args":    eventA.Arguments,
				},
			},
		},
	})

	expected_primitives := []state.Primitive{
		&state.SetPrimitive{
			Path:  []interface{}{},
			Value: map[string]interface{}{},
		},
		&state.SetPrimitive{
			Path:  eventA.Arguments["path"].([]interface{}),
			Value: eventA.Arguments["value"],
		},
		&state.DeletePrimitive{
			Path: eventB.Arguments["path"].([]interface{}),
		},
	}
	for _, ep := range expected_primitives {
		select {
		case primitive := <-primitives:
			switch ep.(type) {
			case *state.SetPrimitive:
				p, ok := primitive.(*state.SetPrimitive)
				if !ok {
					t.Fatalf("Type coercion - expected SET, got DELETE\n%#v\n%#v", p, ep)
				}
				assert.Equal(t, *ep.(*state.SetPrimitive), *p)
			case *state.DeletePrimitive:
				p, ok := primitive.(*state.DeletePrimitive)
				if !ok {
					t.Fatalf("Type coercion - expected DELETE, got SET\n%#v\n%#v", p, ep)
				}
				assert.Equal(t, *ep.(*state.DeletePrimitive), *p)
			default:
				t.Fatal("Was not any known primitive type, wtf")
			}
		case <-time.After(50 * time.Millisecond):
			t.Fatal("Timed out waiting for primitive")
		}
	}
	if len(primitives) > 0 {
		t.Fatal("Unexpected extra primitives")
	}

	expected_export := map[string]interface{}{
		"items": map[string]interface{}{
			"first": "thing",
		},
	}
	assert.Equal(t, expected_export, spt.Simple[0].Export())
	assert.Equal(t, expected_export, spt.Simple[1].Export())
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

func TestSimpleClient_GetTopic(t *testing.T) {
	topic := "http://example.com/deje/some-doc"
	client := NewSimpleClient(topic, nil)
	assert.Equal(t, topic, client.GetTopic())
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
