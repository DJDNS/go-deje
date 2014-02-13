package model

import (
	"encoding/json"
	"github.com/campadrenalin/go-deje/serial"
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
	E := NewEvent("handler name")

	if !Q1.Eq(Q2) {
		t.Fatal("Q1 should equal Q2")
	}

	Q1.Signatures["ident"] = "sig"
	if Q1.Eq(Q2) {
		t.Fatal("Q1 should not equal Q2")
	}

	if Q1.Eq(E) {
		t.Fatal("Q1 should not equal E")
	}
}

func TestQuorum_ToSerial(t *testing.T) {
	Q := NewQuorum("example")
	Q.Signatures["key"] = "value"
	Q.Signatures["this"] = "that"

	SQ := Q.ToSerial()
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
	SQ := serial.Quorum{
		EventHash:  "example",
		Signatures: make(map[string]string),
	}
	SQ.Signatures["hello"] = "world"

	Q := QuorumFromSerial(SQ)

	if Q.EventHash != SQ.EventHash {
		t.Fatalf("EventHash differs: %v vs %v", Q.EventHash, SQ.EventHash)
	}
	if !reflect.DeepEqual(Q.Signatures, SQ.Signatures) {
		t.Fatalf("Signatures differ: %v vs %v", Q.Signatures, SQ.Signatures)
	}
}
