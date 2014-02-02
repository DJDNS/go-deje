package serial

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
