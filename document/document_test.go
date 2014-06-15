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

// Use this type to instigate failures in the JSON module
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

func setupDocument() (Document, *Event, *Quorum) {
	d := NewDocument()
	d.Topic = "Frolicking"

	ev := d.NewEvent("some handler name")
	ev.Arguments["arg"] = "value"
	ev.ParentHash = "Fooblamoose"
	ev.Register()

	q := d.NewQuorum("some event hash")
	q.Signatures["brian blessed"] = "BRIAAAN BLESSED!"
	q.Register()

	return d, &ev, &q
}

func TestDocument_Serialize_WithStuff(t *testing.T) {
	var buffer bytes.Buffer
	d, ev, q := setupDocument()

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

func TestDocument_Deserialize_Empty(t *testing.T) {
	var buffer bytes.Buffer
	d := NewDocument()
	if err := d.Deserialize(&buffer); err == nil {
		t.Fatal("Deserialization should have failed")
	}
}

func TestDocument_Deserialize_Broken(t *testing.T) {
	var buffer BrokenBuffer
	d := NewDocument()
	if err := d.Deserialize(&buffer); err == nil {
		t.Fatal("Deserialization should have failed")
	}
}

func TestDocument_Deserialize_WrongType(t *testing.T) {
	var buffer bytes.Buffer
	d := NewDocument()

	buffer.WriteString(`[]`)
	if err := d.Deserialize(&buffer); err == nil {
		t.Fatal("Deserialization should have failed")
	}
}

func TestDocument_Deserialize_EmptyObject(t *testing.T) {
	var buffer bytes.Buffer
	d := NewDocument()

	buffer.WriteString(`{}`)
	if err := d.Deserialize(&buffer); err != nil {
		t.Fatal(err)
	}
	comparem(t, "", d.Topic, "Topic not set properly")
	comparem(t, 0, len(d.Events), "Events not set properly")
	comparem(t, 0, len(d.Quorums), "Quorums not set properly")
}

func TestDocument_Deserialize_WithStuff(t *testing.T) {
	var buffer bytes.Buffer
	source, ev, q := setupDocument()
	if err := source.Serialize(&buffer); err != nil {
		t.Fatal(err)
	}

	dest := NewDocument()
	if err := dest.Deserialize(&buffer); err != nil {
		t.Fatal(err)
	}
	comparem(t, source.Topic, dest.Topic, "Topic not set properly")
	comparem(t, len(source.Events), len(dest.Events),
		"Wrong number of events")
	comparem(t, len(source.EventsByParent), len(dest.EventsByParent),
		"Did not Register events")
	comparem(t, len(source.Quorums), len(dest.Quorums),
		"Wrong number of quorums")
	comparem(t, len(source.QuorumsByEvent), len(dest.QuorumsByEvent),
		"Did not Register quorums")

	dest_ev := dest.Events[ev.GetKey()]
	dest_q := dest.Quorums[q.GetKey()]
	if !dest_ev.Eq(*ev) {
		t.Fatalf("Events not equal")
	}
	if !dest_q.Eq(*q) {
		t.Fatalf("Quorums not equal")
	}
	comparem(t, &dest, dest_ev.Doc, "Doc pointer not set on Event")
	comparem(t, &dest, dest_q.Doc, "Doc pointer not set on Quorum")
}

func TestDocument_Deserialize_BadKeys(t *testing.T) {
	var buffer bytes.Buffer
	buffer.WriteString(`{"topic":"Frolicking",` +
		`"events":{` +
		`"NotRealKey":{` +
		`"parent":"Fooblamoose","handler":"some handler name",` +
		`"args":{"arg":"value"}` +
		`}},"quorums":{` +
		`"NotRealKey":{` +
		`"event_hash":"some event hash",` +
		`"sigs":{"brian blessed":"BRIAAAN BLESSED!"}` +
		`}}}` +
		"\n",
	)

	dest := NewDocument()
	if err := dest.Deserialize(&buffer); err != nil {
		t.Fatal(err)
	}

	// Check that objects are not present under wrong keys
	if _, ok := dest.Events["NotRealKey"]; ok {
		t.Fatal("Left an Event in under a bad key!")
	}
	if _, ok := dest.Quorums["NotRealKey"]; ok {
		t.Fatal("Left a Quorum in under a bad key!")
	}

	// Check that they are present under the right keys
	_, ev, q := setupDocument()
	if _, ok := dest.Events[ev.GetKey()]; !ok {
		t.Fatal("Event was not registered under correct key")
	}
	if _, ok := dest.Quorums[q.GetKey()]; !ok {
		t.Fatal("Quorum was not registered under correct key")
	}
}
