package model

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestEvent_Serialize(t *testing.T) {
	ev := NewEvent("handler_name")
	ev.Arguments["hello"] = []interface{}{"world", 5}
	ev.Arguments["before"] = nil

	serialized, err := json.Marshal(ev)
	if err != nil {
		t.Fatal("Serialization failed")
	}
	expected := []byte("{" +
		"\"parent\":\"\"," +
		"\"handler\":\"handler_name\"," +
		"\"args\":{" +
		"\"before\":null," +
		"\"hello\":[\"world\",5]" +
		"}" +
		"}")
	if !bytes.Equal(serialized, expected) {
		t.Fatal(string(serialized))
	}
}

func TestEvent_Eq(t *testing.T) {
	A := NewEvent("hello")
	B := NewEvent("hello")
	C := NewEvent("hello")
	D := NewEvent("hello")

	if !(A.Eq(A) && A.Eq(B) && A.Eq(C) && A.Eq(D)) {
		t.Fatal("Freshly initialized events are not equal")
	}

	B.ParentHash = "Ezekiel Wigglesworth"
	if A.Eq(B) {
		t.Fatal("A should not equal B")
	}

	C.HandlerName = "Ezekiel Wigglesworth"
	if A.Eq(B) {
		t.Fatal("A should not equal C")
	}

	D.Arguments["Ezekiel"] = "Wigglesworth"
	if A.Eq(B) {
		t.Fatal("A should not equal D")
	}

	Q := NewQuorum("hello")
	if A.Eq(Q) {
		t.Fatal("A should not equal Q")
	}
}

