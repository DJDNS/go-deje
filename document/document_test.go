package document

import (
	"bytes"
	"errors"
	"testing"
)

func TestNewDocument(t *testing.T) {
	d := NewDocument()

	if d.State.Value == nil {
		t.Fatal("d.State.Value == nil")
	}
	if d.Events == nil {
		t.Fatal("d.Events == nil")
	}
	if d.EventsByParent == nil {
		t.Fatal("d.EventsByParent == nil")
	}
	if d.Quorums == nil {
		t.Fatal("d.Quorums == nil")
	}
	if d.QuorumsByEvent == nil {
		t.Fatal("d.QuorumsByEvent == nil")
	}
}

type BrokenBuffer struct{}

func (bb BrokenBuffer) Read([]byte) (int, error) {
	return 0, errors.New("BrokenBuffer cannot be read >:)")
}
func (bb BrokenBuffer) Write([]byte) (int, error) {
	return 0, errors.New("BrokenBuffer cannot be written to >:)")
}

func TestDocument_Serialize_Empty(t *testing.T) {
	var buffer bytes.Buffer
	d := NewDocument()
	if err := d.Serialize(&buffer); err != nil {
		t.Fatal(err)
	}

	expected := `{"topic":"","events":{},"quorums":{}}` + "\n"
	got := buffer.String()
	if got != expected {
		t.Fatalf("Expected %#v, got %#v", expected, got)
	}
}

func TestDocument_Serialize_Broken(t *testing.T) {
	var buffer BrokenBuffer
	d := NewDocument()
	if err := d.Serialize(&buffer); err == nil {
		t.Fatal("Serialization should fail")
	}
}

func TestDocument_Serialize_WithStuff(t *testing.T) {
	var buffer bytes.Buffer
	d := NewDocument()
	d.Topic = "Frolicking"

	ev := d.NewEvent("some handler name")
	ev.Arguments["arg"] = "value"
	ev.ParentHash = "Fooblamoose"
	ev.Register()

	q := d.NewQuorum("some event hash")
	q.Signatures["brian blessed"] = "BRIAAAN BLESSED!"
	q.Register()

	if err := d.Serialize(&buffer); err != nil {
		t.Fatal(err)
	}
	expected := `{"topic":"Frolicking",` +
		`"events":{` +
		`"` + ev.GetKey() + `":{` +
		`"parent":"Fooblamoose","handler":"some handler name",` +
		`"args":{"arg":"value"}` +
		`}},"quorums":{` +
		`"` + q.GetKey() + `":{` +
		`"event_hash":"some event hash",` +
		`"sigs":{"brian blessed":"BRIAAAN BLESSED!"}` +
		`}}}` +
		"\n"
	got := buffer.String()
	if got != expected {
		t.Fatalf("Expected %#v\n\nGot %#v", expected, got)
	}
}
