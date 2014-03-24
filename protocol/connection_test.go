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
