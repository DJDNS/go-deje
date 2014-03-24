package protocol

import (
	"github.com/campadrenalin/go-deje/logic"
	"github.com/campadrenalin/go-deje/model"
	"github.com/campadrenalin/go-deje/services"
	"testing"
)

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
