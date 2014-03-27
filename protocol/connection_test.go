package protocol

import (
	"bytes"
	"github.com/campadrenalin/go-deje/logic"
	"github.com/campadrenalin/go-deje/model"
	"github.com/campadrenalin/go-deje/services"
	"log"
	"testing"
	"time"
)

var timeout = time.Millisecond

var location = model.IRCLocation{
	Host:    "example.com",
	Port:    6667,
	Channel: "example",
}

func TestConnection_PublishEvent(t *testing.T) {
	d := logic.NewDocument()
	c := services.NewIRCChannel(location)
	conn := NewConnection(d, c)

	ev := model.NewEvent("handler_name")
	err := conn.PublishEvent(ev)
	if err != nil {
		t.Fatal(err)
	}

	if len(c.Outgoing) != 1 {
		t.Fatalf("Expected 1 item in output, got %d", len(c.Outgoing))
	}
	expected := `deje event {"parent":"","handler":"handler_name","args":{}}`
	produced := <-c.Outgoing
	if expected != produced {
		t.Errorf("Expected output: '%s'", expected)
		t.Errorf("Produced output: '%s'", produced)
	}
}

func TestConnection_PublishEvent_Unserializable(t *testing.T) {
	d := logic.NewDocument()
	c := services.NewIRCChannel(location)
	conn := NewConnection(d, c)

	ev := model.NewEvent("handler_name")
	ev.Arguments["evil"] = make(chan int)
	err := conn.PublishEvent(ev)
	if err == nil {
		t.Fatal("Serializing an object containing a chan should fail")
	}

	if len(c.Outgoing) != 0 {
		t.Fatalf("Expected 0 items in output, got %d", len(c.Outgoing))
	}
}

type ParseEventTest struct {
	Name          string
	Input         string
	Event         model.Event
	ShouldSucceed bool
}

func (test ParseEventTest) Run(t *testing.T) {
	d := logic.NewDocument()
	c := services.NewIRCChannel(location)
	conn := NewConnection(d, c)

	err := conn.onEvent(test.Input)
	was_registered := d.Events.Contains(test.Event)
	if test.ShouldSucceed {
		if err != nil {
			t.Log(test.Name)
			t.Fatal(err)
		}
		if !was_registered {
			t.Log(test.Name)
			t.Error("Event not registered")
			t.Fatalf("%#v", d.Events.GetItems())
		}
	} else {
		if err == nil {
			t.Log(test.Name)
			t.Fatal("Parse should fail, but didn't!")
		}
		if was_registered {
			t.Log(test.Name)
			t.Fatal("Event was registered")
		}
	}
}

func TestConnection_onEvent(t *testing.T) {
	ev := model.NewEvent("handler_name")
	ev.Arguments["hello"] = "world"
	ev.ParentHash = "mainstream television"
	tests := []ParseEventTest{
		ParseEventTest{
			"basic",
			`{"parent":"mainstream television",` +
				`"handler":"some_handler",` +
				`"args":{"hello":"world"}}`,
			model.Event{
				ParentHash:  "mainstream television",
				HandlerName: "some_handler",
				Arguments: map[string]interface{}{
					"hello": "world",
				},
			},
			true,
		},
		ParseEventTest{
			"whitespace",
			` {"parent":"",  "handler":"foo","args":{}}  `,
			model.Event{
				HandlerName: "foo",
				Arguments:   map[string]interface{}{},
			},
			true,
		},
		ParseEventTest{
			"crasher",
			"crasher",
			model.Event{},
			false,
		},
		ParseEventTest{
			"leading protocol crap",
			`deje event {"parent":"","handler":"","args":{}}`,
			model.Event{},
			false,
		},
	}

	for _, test := range tests {
		test.Run(t)
	}
}

func TestConnection_OnRecv_NotProtocol(t *testing.T) {
	d := logic.NewDocument()
	c := services.NewIRCChannel(location)
	conn := NewConnection(d, c)

	err := conn.onRecv("cor blimey")
	if err == nil {
		t.Fatal("onRecv should have failed")
	}
	expected_error := `Not a protocol message: "cor blimey"`
	if err.Error() != expected_error {
		t.Errorf("Expected:\n%s", expected_error)
		t.Fatalf("Got:\n%s", err.Error())
	}

	err = conn.onRecv("deje") // No trailing space!
	if err == nil {
		t.Fatal("onRecv should have failed")
	}
	expected_error = `Not a protocol message: "deje"`
	if err.Error() != expected_error {
		t.Errorf("Expected:\n%s", expected_error)
		t.Fatalf("Got:\n%s", err.Error())
	}
}

func TestConnection_OnRecv_UnknownType(t *testing.T) {
	d := logic.NewDocument()
	c := services.NewIRCChannel(location)
	conn := NewConnection(d, c)

	err := conn.onRecv("deje newfoundland")
	if err == nil {
		t.Fatal("onRecv should have failed")
	}
	expected_error := `Not a valid message type: "newfoundland"`
	if err.Error() != expected_error {
		t.Errorf("Expected:\n%s", expected_error)
		t.Fatalf("Got:\n%s", err.Error())
	}
}

func TestConnection_OnRecv_Event(t *testing.T) {
	d := logic.NewDocument()
	c := services.NewIRCChannel(location)
	conn := NewConnection(d, c)

	err := conn.onRecv(`deje event {"parent":"mogwai"}`)
	if err != nil {
		t.Fatal(err)
	}
	ev := model.Event{
		ParentHash: "mogwai", // TODO: Handle missing arguments better
	}
	if !d.Events.Contains(ev) {
		t.Error(d.Events.GetItems())
		t.Fatal("Should have registered event")
	}
}

func testLogger() (*log.Logger, *bytes.Buffer) {
	var b bytes.Buffer
	logger := log.New(&b, "", 0)
	return logger, &b
}

func TestConnection_Run(t *testing.T) {
	d := logic.NewDocument()
	c := services.NewIRCChannel(location)
	conn := NewConnection(d, c)
	logger, buf := testLogger()

	go conn.Run(logger)
	defer conn.Stop()

	ev_string := `deje event {"parent":"","handler":"foo","args":{}}`
	ev := model.Event{
		HandlerName: "foo",
		Arguments:   map[string]interface{}{},
	}

	// Invalid data
	conn.Channel.Incoming <- "invalid data"
	<-time.After(timeout)
	expected_log_output := `Not a protocol message: "invalid data"` + "\n"
	log_output := buf.String()
	if log_output != expected_log_output {
		t.Errorf("Expected:\n'%s'", expected_log_output)
		t.Fatalf("Logged:\n'%s'", log_output)
	}
	num_events := d.Events.Length()
	if num_events != 0 {
		t.Fatalf("Expected 0 registered events, got %d", num_events)
	}

	// Good data
	conn.Channel.Incoming <- ev_string
	<-time.After(timeout)
	log_output = buf.String() // Check that no more errors reported
	if log_output != expected_log_output {
		t.Errorf("Expected:\n'%s'", expected_log_output)
		t.Fatalf("Logged:\n'%s'", log_output)
	}
	num_events = d.Events.Length()
	if num_events != 1 {
		t.Fatalf("Expected 1 registered events, got %d", num_events)
	}
	if !d.Events.Contains(ev) {
		t.Fatal("Should have registered event")
	}
}

func TestConnection_Stop(t *testing.T) {
	d := logic.NewDocument()
	c := services.NewIRCChannel(location)
	conn := NewConnection(d, c)
	logger, buf := testLogger()

	go conn.Run(logger)
	<-time.After(timeout)
	conn.Stop()

	log_output := buf.String()
	expected := "Exiting protocol connection loop\n"
	if log_output != expected {
		t.Errorf("Expected:\n'%s'", expected)
		t.Fatalf("Logged:\n'%s'", log_output)
	}
}
