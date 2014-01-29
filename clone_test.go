package deje

import (
	"testing"
)

func TestCloneMarshal(t *testing.T) {
	m := make(JSONObject)
	m["host"] = "some string"
	m["port"] = 9001
	m["channel"] = "go-nuts"

	loc := new(IRCLocation)
	err := CloneMarshal(m, loc)
	if err != nil {
		t.Fatal("Error in CloneMarshal: %v", err)
	}

	expected := IRCLocation{
		Host:    "some string",
		Port:    9001,
		Channel: "go-nuts",
	}
	if *loc != expected {
		t.Fatalf("Expected %v, got %v", expected, loc)
	}
}

func TestCloneMarshalBadData(t *testing.T) {
	m := make(JSONObject)
	m["ghost"] = "Whatever"
	loc := new(IRCLocation)
	err := CloneMarshal(m, loc)
	if err != nil {
		t.Fatal("CloneMarshal got picky about extra/missing data")
	}

	m["host"] = 5
	err = CloneMarshal(m, loc)
	if err == nil {
		t.Fatal("CloneMarshal should have failed, but didn't")
	}
}
