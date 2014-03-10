package broadcast

import (
	"github.com/eapache/channels"
	"sync"
)

// Represents one of the channels that the Broadcaster sends
// to. It has unbounded buffer capacity (at least until your
// process incurs the wrath of the OOM killer). Even so, you
// really should continue to process your subscription output
// at top speed and keep a clean buffer.
type Subscription struct {
	*channels.InfiniteChannel
	key    int
	source *Broadcaster
}

// Attempt to send some data to a Subscription, and
// return whether it was successful.
func (sub Subscription) Send(data interface{}) bool {
	success := true
	channel := sub.In()
	defer func() {
		if r := recover(); r != nil {
			success = false
		}
	}()
	channel <- data
	return success
}

func (sub Subscription) remove() {
	delete(sub.source.subscriptions, sub.key)
}

// Close a subscription and immediately remove it from the
// Broadcaster's list of subscriptions.
//
// If you close a Subscription's underlying input channel
// directly, the Subscription will be unsubscribed on the
// next broadcast loop. This is not recommended, though,
// because it's an ugly hack that leaves Subscriptions hanging
// around until the next broadcast, and it relies on catching
// panics on channel send. Truly a hideous feature.
func (sub Subscription) Close() {
	b := sub.source
	b.mutex.Lock()
	defer b.mutex.Unlock()

	sub.remove()
	sub.InfiniteChannel.Close()
}

// Ever wanted to send a message to a whole bunch of channels,
// without blocking, in case any of them are doing something
// slow or idiotic? Well now, you can!
//
// Just use our new, handy-dandy Broadcaster struct, and within
// no time, you'll be multicasting messages to a diverse set of
// subscribers, without having to worry about slowing down your
// writing goroutine.
type Broadcaster struct {
	input         chan interface{}
	subscriptions map[int]Subscription
	max_key       int
	mutex         *sync.Mutex
}

// Create a new Broadcaster object, with proper starting parameters
// and an underlying goroutine running behind the scenes.
func NewBroadcaster() *Broadcaster {
	b := &Broadcaster{
		input:         make(chan interface{}),
		subscriptions: make(map[int]Subscription),
		max_key:       0,
		mutex:         new(sync.Mutex),
	}
	go b.run()
	return b
}

// Get input channel, for broadcasting values to all subscribers.
func (b *Broadcaster) In() chan<- interface{} {
	return b.input
}

// Create a new Subscription object, which will recieve broadcasts.
func (b *Broadcaster) Subscribe() Subscription {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	key := b.max_key
	b.max_key = key + 1
	sub := Subscription{
		channels.NewInfiniteChannel(),
		0,
		b,
	}
	b.subscriptions[key] = sub
	return sub
}

// Close input channel, thus shutting down internal goroutine
func (b *Broadcaster) Close() {
	close(b.input)
}

func (b *Broadcaster) run() {
	for {
		data, ok := <-b.input
		if !ok {
			break
		}

		b.mutex.Lock()
		for _, sub := range b.subscriptions {
			open := sub.Send(data)
			if !open {
				sub.remove()
			}
		}
		b.mutex.Unlock()
	}
}
