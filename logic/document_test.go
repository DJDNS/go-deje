package logic

import (
	"github.com/campadrenalin/go-deje/model"
	//"github.com/campadrenalin/go-deje/util"
	"testing"
)

func TestNewDocument(t *testing.T) {
	d := NewDocument()

	if d.Events.GetItems() == nil {
		t.Fatal("d.Events.GetItems() == nil")
	}
	if d.Quorums.GetItems() == nil {
		t.Fatal("d.Quorums.GetItems() == nil")
	}
	if d.Timestamps.GetItems() == nil {
		t.Fatal("d.Timestamps.GetItems() == nil")
	}
}

func TestFromFile(t *testing.T) {
	d := NewDocument()
	df := model.NewDocumentFile()

	df.Channel.Host = "some host"
	df.Channel.Port = 5555 // Interstella?
	df.Channel.Channel = "some channel"

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

	if d.Channel != df.Channel {
		t.Fatal("Channels differ")
	}

	_, ok := d.Events.GetByKey("hello")
	if ok {
		t.Fatal("Event added under the wrong key")
	}

	if d.Events.Length() != 1 {
		t.Fatalf("Wrong num events - expected 1, got %d", d.Events.Length())
	}
	/*
		ev_from_s := EventFromSerial(ev)
		ev_from_d, ok := d.Events.GetByKey(ev_from_s.GetKey())
		if !ok {
			t.Fatal("Could not get event from Document")
		}
		if !ev_from_d.(Event).Eq(ev_from_s) {
			t.Fatalf("%v != %v", ev_from_d, ev_from_s)
		}
	*/

	if d.Quorums.Length() != 1 {
		t.Fatalf("Wrong num quorum - expected 1, got %d", d.Quorums.Length())
	}
	/*
		q_from_s := QuorumFromSerial(q)
		q_from_d, ok := d.Quorums.GetByKey(q_from_s.GetKey())
		if !ok {
			t.Fatal("Could not get quorum from Document")
		}
		if !q_from_d.(Quorum).Eq(q_from_s) {
			t.Fatalf("%v != %v", q_from_d, q_from_s)
		}
	*/
}

func TestToFile(t *testing.T) {
	d := NewDocument()

	d.Channel = model.IRCLocation{
		Host:    "some host",
		Port:    5555,
		Channel: "some channel",
	}

	ev := d.NewEvent("handler name")
	ev.Arguments["hello"] = "world"
	d.Events.Register(ev.Event)

	q := d.NewQuorum("evhash")
	q.Signatures["x"] = "y"
	d.Quorums.Register(q.Quorum)

	df := d.ToFile()

	if df.Channel != d.Channel {
		t.Fatal("Channels differ")
	}

	if d.Events.Length() != 1 {
		t.Fatal("Event conversion failure - wrong num events")
	}

	/*
		ev_to_s := ev.ToSerial()
		ev_df := df.Events[ev.GetKey()]
		hash1, _ := util.HashObject(ev_to_s)
		hash2, _ := util.HashObject(ev_df)
		if hash1 != hash2 {
			t.Fatalf("%v != %v", ev_to_s, ev_df)
		}
	*/

	if d.Quorums.Length() != 1 {
		t.Fatal("Quorum conversion failure - wrong num quorums")
	}

	/*
		q_to_s := q.ToSerial()
		q_df := df.Quorums[q.GetKey()]
		hash1, _ = util.HashObject(q_to_s)
		hash2, _ = util.HashObject(q_df)
		if hash1 != hash2 {
			t.Fatalf("%v != %v", q_to_s, q_df)
		}
	*/
}
