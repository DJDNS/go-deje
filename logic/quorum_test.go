package logic

/*
import (
	"encoding/json"
	"github.com/campadrenalin/go-deje/model"
	"reflect"
	"testing"
)

func TestQuorum_ToSerial(t *testing.T) {
	d := NewDocument()
	Q := d.NewQuorum("example")
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
	SQ := model.Quorum{
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
*/
