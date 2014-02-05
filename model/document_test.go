package deje

import (
	"github.com/campadrenalin/go-deje/serial"
	"github.com/campadrenalin/go-deje/util"
	"testing"
)

func TestNewDocument(t *testing.T) {
	d := NewDocument()

	if d.Events.by_key == nil {
		t.Fatal("d.Events == nil")
	}
	if d.Quorums == nil {
		t.Fatal("d.Quorums == nil")
	}
	if d.Timestamps.Stamps == nil {
		t.Fatal("d.Timestamps.Stamps == nil")
	}
	if d.Timestamps.PerBlock == nil {
		t.Fatal("d.Timestamps.PerBlock == nil")
	}
}

func TestFromFile(t *testing.T) {
	d := NewDocument()
	df := serial.NewDocumentFile()

	df.Channel.Host = "some host"
	df.Channel.Port = 5555 // Interstella?
	df.Channel.Channel = "some channel"

	ev := serial.NewEvent("handler name")
	ev.Arguments["hello"] = "world"
	df.Events["example"] = ev

	d.FromFile(&df)

	if d.Channel != df.Channel {
		t.Fatal("Channels differ")
	}

	if d.Events.Length() != 1 {
		t.Fatal("Event conversion failure - wrong num events")
	}
	ev_from_s := EventFromSerial(ev)
	_, ok := d.Events.GetByKey("hello")
	if ok {
		t.Fatal("Event added under the wrong key")
	}
	ev_from_d, ok := d.Events.GetByKey(ev_from_s.Hash())
	if !ok {
		t.Fatal("Could not get event from Document")
	}
	if ev_from_d.(Event).Hash() != ev_from_s.Hash() {
		t.Fatalf("%v != %v", ev_from_d, ev_from_s)
	}
}

func TestToFile(t *testing.T) {
	d := NewDocument()

	d.Channel = serial.IRCLocation{
		Host:    "some host",
		Port:    5555,
		Channel: "some channel",
	}

	ev := NewEvent("handler name")
	ev.Arguments["hello"] = "world"
	d.Events.Register(ev)

	df := d.ToFile()

	if df.Channel != d.Channel {
		t.Fatal("Channels differ")
	}

	if d.Events.Length() != 1 {
		t.Fatal("Event conversion failure - wrong num events")
	}

	ev_to_s := ev.ToSerial()
	ev_df := df.Events[ev.GetKey()]
	hash1, _ := util.HashObject(ev_to_s)
	hash2, _ := util.HashObject(ev_df)
	if hash1 != hash2 {
		t.Fatalf("%v != %v", ev_to_s, ev_df)
	}
}
