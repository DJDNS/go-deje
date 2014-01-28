package deje

import "testing"

func TestFillStruct(t *testing.T) {
	m := make(JSONObject)
	m["host"] = "some string"
	m["port"] = 9001
	m["channel"] = "go-nuts"

	loc := new(IRCLocation)
	err := FillStruct(m, loc)
	if err != nil {
		t.Fatal("Error in FillStruct: %v", err)
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

func TestFillStructBadData(t *testing.T) {
	m := make(JSONObject)
	m["ghost"] = "Whatever"
	loc := new(IRCLocation)
	err := FillStruct(m, loc)
	if err != nil {
		t.Fatal("FillStruct got picky about extra/missing data")
	}

	m["host"] = 5
	err = FillStruct(m, loc)
	if err == nil {
		t.Fatal("FillStruct should have failed, but didn't")
	}
}

func TestDSGetChannel(t *testing.T) {
	ds := make(DocumentState)

	channel := make(JSONObject)
	channel["host"] = "some string"
	channel["port"] = 9001
	channel["channel"] = "go-nuts"
	ds["channel"] = channel

	loc, err := ds.GetChannel()
	if err != nil {
		t.Fatal(err)
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

func TestDSGetChannelBadData(t *testing.T) {
	ds := make(DocumentState)

	_, err := ds.GetChannel()
	if err == nil {
		t.Fatal("GetChannel should have failed, but didn't")
	}

	ds["channel"] = 4
	_, err = ds.GetChannel()
	if err == nil {
		t.Fatal("GetChannel should have failed, but didn't")
	}

	channel := make(JSONObject)
	channel["port"] = "string port"
	ds["channel"] = channel
	_, err = ds.GetChannel()
	if err == nil {
		t.Fatal("GetChannel should have failed, but didn't")
	}
}
