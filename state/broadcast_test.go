package state

import (
	"reflect"
	"testing"
	"time"
)

var timeout = time.Millisecond

func TestNewPrimitiveBroadcaster(t *testing.T) {
	b := NewPrimitiveBroadcaster()
	if b.Broadcaster == nil {
		t.Fatal("b.Broadcaster not initialized")
	}
}

func TestPrimitiveSubscription_Out(t *testing.T) {
	b := NewPrimitiveBroadcaster()
	sub := b.Subscribe()
	data := &SetPrimitive{
		[]interface{}{"hello"},
		"world",
	}

	b.Send(data)
	select {
	case result := <-sub.Out():
		if !reflect.DeepEqual(result, data) {
			t.Fatalf("Expected %#v, got %#v", data, result)
		}
	case <-time.After(timeout):
		t.Fatal("No data passed through")
	}
}

func TestPrimitiveSubscription_Len(t *testing.T) {
	b := NewPrimitiveBroadcaster()
	sub := b.Subscribe()
	data := &SetPrimitive{
		[]interface{}{"hello"},
		"world",
	}
	num_sends := 30
	num_recvs := 5

	for i := 0; i < num_sends; i++ {
		b.Send(data)
	}
	<-time.After(timeout) // Wait for things to get settled
	length := sub.Len()
	if length != num_sends {
		t.Errorf("Sub len: %d", sub.sub.Len())
		t.Errorf("out len: %d", len(sub.out))
		t.Errorf("sending: %d", sub.sending)
		t.Fatalf("Expected %d, got %d", num_sends, length)
	}

	for i := 0; i < num_recvs; i++ {
		select {
		case result := <-sub.Out():
			if !reflect.DeepEqual(result, data) {
				t.Errorf("Iteration %d:", i)
				t.Fatalf("Expected %#v, got %#v", data, result)
			}
		case <-time.After(timeout):
			t.Fatal("No data passed through")
		}
	}
	<-time.After(timeout) // Wait for things to get settled

	expected_len := num_sends - num_recvs
	length = sub.Len()
	if length != expected_len {
		t.Errorf("Sub len: %d", sub.sub.Len())
		t.Errorf("out len: %d", len(sub.out))
		t.Errorf("sending: %d", sub.sending)
		t.Fatalf("Expected %d, got %d", expected_len, length)
	}
}
