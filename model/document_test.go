package deje

import (
	"github.com/campadrenalin/go-deje/serial"
	"testing"
)

func TestNewDocument(t *testing.T) {
	d := NewDocument()

	if d.Events == nil {
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

	if len(d.Events) != 1 {
		t.Fatal("Event conversion failure - wrong num events")
	}
	ev_from_s := EventFromSerial(ev)
	ev_from_d := d.Events["example"]
	if ev_from_d.Hash() != ev_from_s.Hash() {
		t.Fatalf("%v != %v", ev_from_d, ev_from_s)
	}
}
