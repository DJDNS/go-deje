package deje

import (
    "encoding/json"
    "bytes"
    "strings"
    "testing"
)

func TestEvent_Serialize(t *testing.T) {
    ev := NewEvent(SHA1Hash{}, "handler_name")
    ev.Arguments["hello"] = []interface{}{"world", 5}
    ev.Arguments["before"] = nil

    serialized, err := json.Marshal(ev)
    if err != nil {
        t.Fatal("Serialization failed")
    }
    expected := []byte(strings.Join([]string{
        "{\"parent\":[" + strings.Repeat("0,", 19) + "0],",
        "\"handler\":\"handler_name\",",
        "\"args\":{",
            "\"before\":null,",
            "\"hello\":[\"world\",5]",
        "}}",
    }, ""))
    if !bytes.Equal(serialized, expected) {
        t.Fatal(string(serialized))
    }
}

func TestEventSet_GetRoot_NoElements(t *testing.T) {
    set := make(EventSet)
    _, ok := set.GetRoot(NewEvent(SHA1Hash{1,2,3}, "handler_name"))
    if ok {
        t.Fatal("GetRoot should have failed, but returned ok == true")
    }
}

func TestEventSet_GetRoot(t *testing.T) {
    set := make(EventSet)
    first := NewEvent(SHA1Hash{}, "first")
    second := NewEvent(first.Hash(), "second")
    third := NewEvent(second.Hash(), "third")

    events := []Event{first, second, third}
    for _, ev := range events {
        set.Register(ev)
    }

    for _, ev := range events {
        found, ok := set.GetRoot(ev)
        if ! ok {
            t.Fatal("GetRoot failed")
        }
        if found.HandlerName != "first" {
            t.Fatal("Did not get correct event")
        }
    }
}
