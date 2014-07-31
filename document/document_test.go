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

	expected := `{"topic":"","events":{},"quorums":{},"timestamps":[]}` + "\n"
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

func setupDocument() (Document, []*Event, []*Quorum) {
	d := NewDocument()
	d.Topic = "Frolicking"

	// These values have been adjusted to ensure that slice
	// position reflects hash order
	events := make([]*Event, 2)
	ev0 := d.NewEvent("some handler name")
	ev0.Arguments["arg"] = "value"
	ev0.ParentHash = "Fooblamoose"
	ev0.Register()
	events[0] = &ev0

	ev1 := d.NewEvent("other handler name")
	ev1.Arguments["cow"] = "moo"
	ev1.Register()
	events[1] = &ev1

	// Same as with events - adjusted to maintain that
	// hash order == slice order.
	quorums := make([]*Quorum, 2)
	q0 := d.NewQuorum("some event hash")
	q0.Signatures["brian blessed"] = "BRIAAAN BLESSED!"
	q0.Register()
	quorums[0] = &q0

	q1 := d.NewQuorum("other event hash")
	q1.Signatures["John Hancock"] = "<swoopy cursive>"
	q1.Register()
	quorums[1] = &q1

	return d, events, quorums
}

func TestDocument_Serialize_WithStuff(t *testing.T) {
	var buffer bytes.Buffer
	d, ev, q := setupDocument()

	if err := d.Serialize(&buffer); err != nil {
		t.Fatal(err)
	}
	expected := `{"topic":"Frolicking",` +
		`"events":{` +
		`"` + ev[0].GetKey() + `":{` +
		`"parent":"Fooblamoose","handler":"some handler name",` +
		`"args":{"arg":"value"}` +
		`},"` + ev[1].GetKey() + `":{` +
		`"parent":"","handler":"other handler name",` +
		`"args":{"cow":"moo"}` +
		`}},"quorums":{` +
		`"` + q[0].GetKey() + `":{` +
		`"event_hash":"some event hash",` +
		`"sigs":{"brian blessed":"BRIAAAN BLESSED!"}` +
		`},"` + q[1].GetKey() + `":{` +
		`"event_hash":"other event hash",` +
		`"sigs":{"John Hancock":"\u003cswoopy cursive\u003e"}` +
		`}},"timestamps":[]}` +
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

	for i := range ev {
		dest_ev := dest.Events[ev[i].GetKey()]
		if !dest_ev.Eq(*ev[i]) {
			t.Fatalf("Events not equal: %s", ev[i].HandlerName)
		}
		comparem(t, &dest, dest_ev.Doc, "Doc pointer not set on Event")
	}
	for i := range q {
		dest_q := dest.Quorums[q[i].GetKey()]
		if !dest_q.Eq(*q[i]) {
			t.Fatalf("Quorums not equal: %d", i)
		}
		comparem(t, &dest, dest_q.Doc, "Doc pointer not set on Quorum")
	}
}

func TestDocument_Deserialize_BadKeys(t *testing.T) {
	var buffer bytes.Buffer
	buffer.WriteString(`{"topic":"Frolicking",` +
		`"events":{` +
		`"NotRealKey":{` +
		`"parent":"Fooblamoose","handler":"some handler name",` +
		`"args":{"arg":"value"}` +
		`},"AlsoNotReal":{` +
		`"parent":"","handler":"some other handler name",` +
		`"args":{}` +
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
	if _, ok := dest.Events[ev[0].GetKey()]; !ok {
		t.Fatal("Event was not registered under correct key")
	}
	if _, ok := dest.Quorums[q[0].GetKey()]; !ok {
		t.Fatal("Quorum was not registered under correct key")
	}
}
