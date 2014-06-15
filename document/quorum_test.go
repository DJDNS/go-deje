package document

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestQuorum_GetKey(t *testing.T) {
	A := NewQuorum("hello")
	B := NewQuorum("hello")

	if A.GetKey() != B.GetKey() {
		t.Fatal("A should equal B")
	}

	A.Signatures["ident"] = "sig"
	if A.GetKey() == B.GetKey() {
		t.Fatal("A should not equal B")
	}
}

func TestQuorum_GetGroupKey(t *testing.T) {
	A := NewQuorum("hello")
	B := NewQuorum("hello")

	if A.GetGroupKey() != B.GetGroupKey() {
		t.Fatal("A should equal B")
	}

	// Signature content should not change group key
	A.Signatures["ident"] = "sig"
	if A.GetGroupKey() != B.GetGroupKey() {
		t.Fatal("A should equal B")
	}

	// Different EventHash should alter group key
	A.EventHash = "world"
	if A.GetGroupKey() == B.GetGroupKey() {
		t.Fatal("A should not equal B")
	}
}

func TestQuorum_Eq(t *testing.T) {
	Q1 := NewQuorum("hello")
	Q2 := NewQuorum("hello")

	if !Q1.Eq(Q2) {
		t.Fatal("Q1 should equal Q2")
	}

	Q1.Signatures["ident"] = "sig"
	if Q1.Eq(Q2) {
		t.Fatal("Q1 should not equal Q2")
	}
}

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

	expected := "{" +
		"\"event_hash\":\"example\"," +
		"\"sigs\":{" +
		"\"key\":\"value\"," +
		"\"this\":\"that\"" +
		"}}"
	got, err := json.Marshal(Q)
	if err != nil {
		t.Fatal("Serialization failed", err)
	}
	gotstr := string(got)
	if gotstr != expected {
		t.Fatalf("Expected %v, got %v", expected, gotstr)
	}
}

func TestQuorumFromSerial(t *testing.T) {
	source := "{" +
		"\"event_hash\":\"example\"," +
		"\"sigs\":{" +
		"\"hello\":\"world\"" +
		"}}"
	var got Quorum
	err := json.Unmarshal([]byte(source), &got)
	if err != nil {
		t.Fatal("Deserialization failed:", err)
	}

	expected := Quorum{
		EventHash: "example",
		Signatures: map[string]string{
			"hello": "world",
		},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Quorums differ: %v vs %v", got, expected)
	}
}
