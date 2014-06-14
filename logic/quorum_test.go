package logic

import (
	"encoding/json"
	"github.com/campadrenalin/go-deje/model"
	"reflect"
	"testing"
)

func TestQuorum_Register(t *testing.T) {
	d := NewDocument()
	q_sweet := d.NewQuorum("sweet")
	q_saltyA := d.NewQuorum("salty")
	q_saltyB := d.NewQuorum("salty")

	q_saltyA.Signatures["Pretzels"] = "pretzy montgomery"
	q_saltyB.Signatures["Fritos"] = "frito bandito"

	quorums := []*Quorum{&q_sweet, &q_saltyA, &q_saltyB}
	for _, q := range quorums {
		q.Register()
	}

	// Test that main set registered correctly
	expected_quorums := QuorumSet{
		q_sweet.GetKey():  &q_sweet,
		q_saltyA.GetKey(): &q_saltyA,
		q_saltyB.GetKey(): &q_saltyB,
	}
	if !reflect.DeepEqual(d.Quorums, expected_quorums) {
		t.Fatalf("Expected %#v\nGot %#v", expected_quorums, d.Quorums)
	}

	// Test that groupings registered correctly
	expected_groups := map[string]QuorumSet{
		"sweet": QuorumSet{
			q_sweet.GetKey(): &q_sweet,
		},
		"salty": QuorumSet{
			q_saltyA.GetKey(): &q_saltyA,
			q_saltyB.GetKey(): &q_saltyB,
		},
	}
	if !reflect.DeepEqual(d.QuorumsByEvent, expected_groups) {
		t.Fatalf("Expected %#v\nGot %#v", expected_groups, d.QuorumsByEvent)
	}
}

func TestQuorum_Unregister(t *testing.T) {
	d := NewDocument()
	q_sweet := d.NewQuorum("sweet")
	q_saltyA := d.NewQuorum("salty")
	q_saltyB := d.NewQuorum("salty")

	q_saltyA.Signatures["Pretzels"] = "pretzy montgomery"
	q_saltyB.Signatures["Fritos"] = "frito bandito"

	quorums := []*Quorum{&q_sweet, &q_saltyA, &q_saltyB}
	for _, q := range quorums {
		q.Register()
	}

	// Unregister from multi-element group
	q_saltyA.Unregister()
	expected_quorums := QuorumSet{
		q_sweet.GetKey():  &q_sweet,
		q_saltyB.GetKey(): &q_saltyB,
	}
	if !reflect.DeepEqual(d.Quorums, expected_quorums) {
		t.Fatalf("Expected %#v\nGot %#v", expected_quorums, d.Quorums)
	}
	expected_groups := map[string]QuorumSet{
		"sweet": QuorumSet{
			q_sweet.GetKey(): &q_sweet,
		},
		"salty": QuorumSet{
			q_saltyB.GetKey(): &q_saltyB,
		},
	}
	if !reflect.DeepEqual(d.QuorumsByEvent, expected_groups) {
		t.Fatalf("Expected %#v\nGot %#v", expected_groups, d.QuorumsByEvent)
	}

	// Unregister from single-element group
	q_sweet.Unregister()
	expected_quorums = QuorumSet{
		q_saltyB.GetKey(): &q_saltyB,
	}
	if !reflect.DeepEqual(d.Quorums, expected_quorums) {
		t.Fatalf("Expected %#v\nGot %#v", expected_quorums, d.Quorums)
	}
	expected_groups = map[string]QuorumSet{
		"salty": QuorumSet{
			q_saltyB.GetKey(): &q_saltyB,
		},
	}
	if !reflect.DeepEqual(d.QuorumsByEvent, expected_groups) {
		t.Fatalf("Expected %#v\nGot %#v", expected_groups, d.QuorumsByEvent)
	}

	// Make sure that Unregistering multiple times is okay
	q_sweet.Unregister()
	if !reflect.DeepEqual(d.Quorums, expected_quorums) {
		t.Fatalf("Expected %#v\nGot %#v", expected_quorums, d.Quorums)
	}
	if !reflect.DeepEqual(d.QuorumsByEvent, expected_groups) {
		t.Fatalf("Expected %#v\nGot %#v", expected_groups, d.QuorumsByEvent)
	}
}

func TestQuorum_ToSerial(t *testing.T) {
	d := NewDocument()
	Q := d.NewQuorum("example")
	Q.Signatures["key"] = "value"
	Q.Signatures["this"] = "that"

	SQ := Q.Quorum
	expected := "{" +
		"\"event_hash\":\"example\"," +
		"\"sigs\":{" +
		"\"key\":\"value\"," +
		"\"this\":\"that\"" +
		"}}"
	got, err := json.Marshal(SQ)
	if err != nil {
		t.Fatal("Serialization failed", err)
	}
	gotstr := string(got)
	if gotstr != expected {
		t.Fatalf("Expected %v, got %v", expected, gotstr)
	}
}

func TestQuorumFromSerial(t *testing.T) {
	d := NewDocument()
	SQ := model.Quorum{
		EventHash:  "example",
		Signatures: make(map[string]string),
	}
	SQ.Signatures["hello"] = "world"

	Q := Quorum{SQ, &d}

	if Q.EventHash != SQ.EventHash {
		t.Fatalf("EventHash differs: %v vs %v", Q.EventHash, SQ.EventHash)
	}
	if !reflect.DeepEqual(Q.Signatures, SQ.Signatures) {
		t.Fatalf("Signatures differ: %v vs %v", Q.Signatures, SQ.Signatures)
	}
}
