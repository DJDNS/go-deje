package util

import (
	"encoding/json"
	"testing"
)

func TestHashObjectBasic(t *testing.T) {
	obj := make(jsonObject)
	obj["x"] = "y"
	obj["z"] = []interface{}{8, 9, nil, true}

	// For debugging
	marshalled, err := json.Marshal(obj)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Marshalled to '%v'", string(marshalled))
	}

	// Obtained with:
	// echo -n '{"x":"y","z":[8,9,null,true]}' | sha1sum
	expected := "b39d52797d2e72ddbe4f2b940a6700122d288a0c"

	hash, err := HashObject(obj)
	if err != nil {
		t.Fatal(err)
	}

	if hash != expected {
		t.Fatalf("Expected %v, got %v", expected, hash)
	}
}

func TestHashObjectEmpty(t *testing.T) {
	loc := new(ircLocation)
	hash, err := HashObject(loc)
	if err != nil {
		t.Fatal(err)
	}

	// For debugging
	marshalled, err := json.Marshal(loc)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Marshalled to '%v'", string(marshalled))
	}

	// Obtained with:
	// echo -n '{"host":"","port":0,"channel":""}' | sha1sum
	expected := "47f58ae5a60f9e23a2c90e02faf4040ac6b0a98b"

	if hash != expected {
		t.Fatalf("Expected %v, got %v", expected, hash)
	}
}

func TestHashObjectPopulated(t *testing.T) {
	loc := new(ircLocation)
	loc.Host = "example.com"
	loc.Port = 666
	loc.Channel = "mtv"

	hash, err := HashObject(loc)
	if err != nil {
		t.Fatal(err)
	}

	// For debugging
	marshalled, err := json.Marshal(loc)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Marshalled to '%v'", string(marshalled))
	}

	// Obtained with:
	// echo -n '{"host":"example.com","port":666,"channel":"mtv"}' | sha1sum
	expected := "6f226d7455fbb88772b8e933009aa8e2bf7800df"

	if hash != expected {
		t.Fatalf("Expected %v, got %v", expected, hash)
	}
}

func TestHashObjectUnmarshallable(t *testing.T) {
	c := make(chan int)
	_, err := HashObject(c)
	if err == nil {
		t.Fatal("HashObject should have choked on unmarshallable object")
	}
}
