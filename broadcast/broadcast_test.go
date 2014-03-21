package broadcast

import (
	"reflect"
	"testing"
	"time"
)

var timeout = time.Millisecond

func TestNewBroadcaster(t *testing.T) {
	b := NewBroadcaster()

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

	b.Send("hello")
}

func TestBroadcaster_Send_OneSub(t *testing.T) {
	b := NewBroadcaster()
	defer b.Close()

	data := "example"
	sub := b.Subscribe()
	b.Send(data)

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
	subs := make([]*Subscription, 10)
	for i, _ := range subs {
		subs[i] = b.Subscribe()
	}

	b.Send(data)

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
	for i := 0; i < max_data; i++ {
		b.Send(i)
	}

	length := sub.Len()
	if length != max_data {
		t.Fatalf("Expected length %d, got %d", max_data, length)
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

func TestBroadcaster_Overflow(t *testing.T) {
	b := NewBroadcaster()
	defer b.Close()

	// Use a flood of data to exceed capacity
	sub := b.Subscribe()
	out := sub.Out()
	capacity := cap(out)

	// Fill exactly to capacity
	for i := 0; i < capacity; i++ {
		b.Send(i)
	}
	if _, ok := <-out; !ok {
		t.Fatal("Chan closed after filling to capacity")
	}
	if sub.Overflowed {
		t.Fatal("sub.Overflowed after filling to capacity")
	}

	// Need to send two to overflow - consumed one earlier
	b.Send(capacity + 1)
	b.Send(capacity + 2)
	if !sub.Overflowed {
		t.Fatal("!sub.Overflowed after filling past capacity")
	}
	if sub.Len() != capacity {
		t.Fatalf("Expected len %d, got %d", capacity, sub.Len())
	}

	// Eat through queue to prove closed-ness
	for i := 0; i < capacity; i++ {
		select {
		case <-out:
		default:
			t.Fatalf("Could not get item %d", i)
		}
	}
	if _, ok := <-out; ok {
		t.Fatal("Chan open after filling past capacity")
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
		b.Send("some data")
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

func TestBroadcaster_UnsubscribeCorrectSub(t *testing.T) {
	b := NewBroadcaster()
	defer b.Close()

	happy := b.Subscribe()
	doomed := b.Subscribe()
	indifferent := b.Subscribe()
	doomed.Close()

	b.mutex.Lock()
	defer b.mutex.Unlock()

	expected := map[int]*Subscription{
		happy.key:       happy,
		indifferent.key: indifferent,
	}
	if !reflect.DeepEqual(b.subscriptions, expected) {
		t.Errorf("Expected %#v", expected)
		t.Fatalf("Got %#v", b.subscriptions)
	}
}

func TestBroadcaster_CleansUpSubscriptions(t *testing.T) {
	b := NewBroadcaster()
	_ = b.Subscribe()
	b.Close()

	// Allow some time for loop, make sure we're synced
	<-time.After(timeout)
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if len(b.subscriptions) > 0 {
		t.Fatal("Still has subscriptions after closing!")
	}
}

func TestBroadcaster_SubscribeAfterClose(t *testing.T) {
	b := NewBroadcaster()
	b.Close()

	paniced := false
	defer func() {
		if r := recover(); r != nil {
			paniced = true
		}
	}()

	_ = b.Subscribe()

	if !paniced {
		t.Fatal("Should panic on Subscribe-after-Close")
	}
}
