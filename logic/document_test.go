package logic

import (
	"github.com/campadrenalin/go-deje/model"
	"github.com/campadrenalin/go-deje/util"
	"testing"
)

func TestNewDocument(t *testing.T) {
	d := NewDocument()

	if d.State.Value == nil {
		t.Fatal("d.State.Value == nil")
	}
	if d.Events.GetItems() == nil {
		t.Fatal("d.Events.GetItems() == nil")
	}
	if d.Quorums.GetItems() == nil {
		t.Fatal("d.Quorums.GetItems() == nil")
	}
}

func TestFromFile(t *testing.T) {
	d := NewDocument()
	df := model.NewDocumentFile()

	df.Topic = "com.example.deje.5555"

	ev := d.NewEvent("handler name")
	ev.Arguments["hello"] = "world"
	df.Events["example"] = ev.Event

	q := model.Quorum{
		EventHash:  "evhash",
		Signatures: make(map[string]string),
	}
	q.Signatures["x"] = "y"
	df.Quorums["example"] = q

	d.FromFile(&df)

	if d.Topic != df.Topic {
		t.Fatal("Topics differ")
	}

	_, ok := d.Events.GetByKey("hello")
	if ok {
		t.Fatal("Event added under the wrong key")
	}

	if d.Events.Length() != 1 {
		t.Fatalf("Wrong num events - expected 1, got %d", d.Events.Length())
	}
	ev_from_s := ev.Event
	ev_from_d, ok := d.Events.GetByKey(ev_from_s.GetKey())
	if !ok {
		t.Fatal("Could not get event from Document")
	}
	if !ev_from_d.(model.Event).Eq(ev_from_s) {
		t.Fatalf("%v != %v", ev_from_d, ev_from_s)
	}

	if d.Quorums.Length() != 1 {
		t.Fatalf("Wrong num quorum - expected 1, got %d", d.Quorums.Length())
	}
	q_from_d, ok := d.Quorums.GetByKey(q.GetKey())
	if !ok {
		t.Fatal("Could not get quorum from Document")
	}
	if !q_from_d.(model.Quorum).Eq(q) {
		t.Fatalf("%v != %v", q_from_d, q)
	}
}

func TestToFile(t *testing.T) {
	d := NewDocument()

	d.Topic = "com.example.deje.5555"

	ev := d.NewEvent("handler name")
	ev.Arguments["hello"] = "world"
	ev.Register()

	q := d.NewQuorum("evhash")
	q.Signatures["x"] = "y"
	q.Register()

	df := d.ToFile()

	if df.Topic != d.Topic {
		t.Fatal("Topics differ")
	}

	if d.Events.Length() != 1 {
		t.Fatal("Event conversion failure - wrong num events")
	}

	ev_to_s := ev.Event
	ev_df := df.Events[ev.GetKey()]
	hash1, _ := util.HashObject(ev_to_s)
	hash2, _ := util.HashObject(ev_df)
	if hash1 != hash2 {
		t.Fatalf("%v != %v", ev_to_s, ev_df)
	}

	if d.Quorums.Length() != 1 {
		t.Fatal("Quorum conversion failure - wrong num quorums")
	}

	q_to_s := q.Quorum
	q_df := df.Quorums[q.GetKey()]
	hash1, _ = util.HashObject(q_to_s)
	hash2, _ = util.HashObject(q_df)
	if hash1 != hash2 {
		t.Fatalf("%v != %v", q_to_s, q_df)
	}
}
