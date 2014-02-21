package services

import (
	"github.com/campadrenalin/go-deje/model"
	"testing"
)

func TestDIS_GetChannel(t *testing.T) {
	location := model.IRCLocation{
		Host:    "x",
		Port:    10,
		Channel: "y",
	}
	dis := DummyIRCService{}
	channel := dis.GetChannel(location)

	if channel.Incoming == nil {
		t.Fatal("channel.Incoming not initialized")
	}
	if channel.Outgoing == nil {
		t.Fatal("channel.Outgoing not initialized")
	}
	if channel.Location != location {
		t.Fatalf(
			"Locations don't match: %#v != %#v",
			channel.Location,
			location,
		)
	}
}
