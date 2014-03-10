package broadcast

import (
	"testing"
	"time"
)

var timeout = time.Millisecond

func TestNewBroadcaster(t *testing.T) {
	b := NewBroadcaster()
	defer b.Close()

	if b.input == nil {
		t.Fatal("b.input not initialized")
	}
	if b.subscriptions == nil {
		t.Fatal("b.subscriptions not initialized")
	}
}

func TestBroadcaster_Subscribe(t *testing.T) {
	b := NewBroadcaster()
	defer b.Close()

	sub := b.Subscribe()

	if sub.Len() != 0 {
		t.Fatal("Sub should be acting like valid, empty InfiniteChannel")
	}
}

func TestBroadcaster_Send_NoSubs(t *testing.T) {
	b := NewBroadcaster()
	defer b.Close()

	input := b.In()

	select {
	case input <- "hello":
	case <-time.After(timeout):
		t.Fatal("Send should not have blocked")
	}
}

func TestBroadcaster_Send_OneSub(t *testing.T) {
	b := NewBroadcaster()
	defer b.Close()

	data := "example"
	input := b.In()
	sub := b.Subscribe()
	input <- data

	select {
	case recvd := <-sub.Out():
		if recvd != data {
			t.Fatal("Expected %v, got %v", data, recvd)
		}
	case <-time.After(timeout):
		t.Fatal("Receive should not have blocked")
	}

}

func TestBroadcaster_Send_MultipleSubs(t *testing.T) {
	b := NewBroadcaster()
	defer b.Close()

	data := "example"
	subs := make([]Subscription, 10)
	for i, _ := range subs {
		subs[i] = b.Subscribe()
	}

	b.In() <- data

	for i, sub := range subs {
		select {
		case recvd := <-sub.Out():
			if recvd != data {
				t.Fatal("Expected %v, got %v", data, recvd)
			}
		case <-time.After(timeout):
			t.Fatalf("Receive should not have blocked (sub %d)", i)
		}
	}
}

func TestBroadcaster_Send_MultipleData(t *testing.T) {
	b := NewBroadcaster()
	defer b.Close()

	// Use a ton of data to stress-test ordering
	max_data := 200
	sub := b.Subscribe()
	input := b.In()
	for i := 0; i < max_data; i++ {
		input <- i
	}

	for i := 0; i < max_data; i++ {
		select {
		case recvd := <-sub.Out():
			if recvd != i {
				t.Fatal("Expected %v, got %v", i, recvd)
			}
		case <-time.After(timeout):
			t.Fatalf("Receive should not have blocked (item %d)", i)
		}
	}
}

func TestBroadcaster_ParallelSubscribe(t *testing.T) {
	b := NewBroadcaster()
	defer b.Close()

	go b.Subscribe()
	go b.Subscribe()

	<-time.After(timeout) // Just long enough to yield to Subscribes
	b.mutex.Lock()        // Make race detector happy
	defer b.mutex.Unlock()

	if len(b.subscriptions) < 2 {
		t.Fatal("Race condition triggered")
	}
}

// Should basically only trigger the race detector
func TestBroadcaster_SubscribeVsBroadcast(t *testing.T) {
	b := NewBroadcaster()
	defer b.Close()

	go b.Subscribe()
	go func() {
		b.In() <- "some data"
	}()

	// Keep broadcaster alive awhile
	<-time.After(timeout)
	b.mutex.Lock()
	defer b.mutex.Unlock()
}

func TestBroadcaster_Unsubscribe(t *testing.T) {
	b := NewBroadcaster()
	defer b.Close()

	sub := b.Subscribe()
	sub.Close()

	if len(b.subscriptions) > 0 {
		t.Fatal("Sub should immediately be removed from list")
	}
}

func TestBroadcaster_UnsubscribeUgly(t *testing.T) {
	b := NewBroadcaster()
	defer b.Close()

	sub := b.Subscribe()
	sub.InfiniteChannel.Close()

	if len(b.subscriptions) != 1 {
		t.Fatal("Shouldn't be immediately removed - b can't know")
	}

	b.In() <- "some data"

	// Allow some time for loop, make sure we're synced
	<-time.After(timeout)
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if len(b.subscriptions) > 0 {
		t.Fatal("Sub should be removed from list on send iteration")
	}
}
