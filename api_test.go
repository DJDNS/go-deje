package deje

import (
	"github.com/campadrenalin/go-deje/model"
	"testing"
)

func TestGetDocument(t *testing.T) {
	con := NewDEJEController()
	loc := model.IRCLocation{"hello", 1234, "world"}
	doc := con.GetDocument(loc)

	if doc.Channel != loc {
		t.Fatalf("Expected %v, got %v", loc, doc.Channel)
	}
}

func TestGetDocument_Persistence(t *testing.T) {
	con := NewDEJEController()
	loc := model.IRCLocation{"hello", 1234, "world"}

	docA := con.GetDocument(loc)
	eventA := model.NewEvent("foo")
	docA.Events.Register(eventA)

	docB := con.GetDocument(loc)
	eventB := model.NewEvent("bar")
	docB.Events.Register(eventB)

	if docA.Channel != docB.Channel {
		t.Fatalf("Locations differ, %v vs %v", docA.Channel, docB.Channel)
	}
	if docA.Events != docB.Events {
		t.Fatalf("Events differ, %v vs %v", docA.Events, docB.Events)
	}
	if docA.Quorums != docB.Quorums {
		t.Fatalf("Quorums differ, %v vs %v", docA.Quorums, docB.Quorums)
	}
	if docA.Timestamps != docB.Timestamps {
		t.Fatalf("Timestamps differ, %v vs %v", docA.Timestamps, docB.Timestamps)
	}

	if !docA.Events.Contains(eventB) {
		t.Fatal("docA should contain eventB")
	}
	if !docB.Events.Contains(eventA) {
		t.Fatal("docB should contain eventA")
	}
}
