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

	if channel.Channel == nil {
		t.Fatal("channel.Channel not initialized")
	}
	if channel.Location != location {
		t.Fatalf(
			"Locations don't match: %#v != %#v",
			channel.Location,
			location,
		)
	}
}
