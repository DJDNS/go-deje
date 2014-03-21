/*
Generic package for single-producer multi-consumer queues.

You should generally wrap this in a type-safe wrapper class
for convenience and safety, unless you actually want to pass
interface{}s around.

Be aware that Subscriptions can "break" if you don't receive
from them fast enough to drain the queue, in which case, it
is your responsibility to resubscribe and recover.
*/
package broadcast

import "sync"

// Represents one of the channels that the Broadcaster sends
// to. If capacity is exhausted, the Subscription will be
// closed, with Overflowed set to true.
type Subscription struct {
	channel    chan interface{}
	key        int
	source     *Broadcaster
	Overflowed bool
}

// Get length of inner channel.
//
// Equivalent to len(sub.Out()).
func (sub *Subscription) Len() int {
	return len(sub.channel)
}

// Get receive-only channel of Subscription.
func (sub *Subscription) Out() <-chan interface{} {
	return sub.channel
}

// Attempt to send some data to a Subscription, and
// return whether it was successful.
func (sub *Subscription) Send(data interface{}) bool {
	channel := sub.channel

	// Non-blocking send - fail if blocked
	select {
	case channel <- data:
		return true
	default:
		return false
	}
}

func (sub *Subscription) remove() {
	delete(sub.source.subscriptions, sub.key)
}

func (sub *Subscription) do_close() {
	sub.remove()
	close(sub.channel)
}

// Close a subscription and immediately remove it from the
// Broadcaster's list of subscriptions.
//
// DO NOT close a Subscription's underlying channel directly.
func (sub *Subscription) Close() {
	b := sub.source
	b.mutex.Lock()
	defer b.mutex.Unlock()

	sub.do_close()
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
	subscriptions map[int]*Subscription
	max_key       int
	mutex         *sync.Mutex
	closed        bool
}

// Create a new Broadcaster object, with proper starting parameters
// and an underlying goroutine running behind the scenes.
func NewBroadcaster() *Broadcaster {
	b := &Broadcaster{
		subscriptions: make(map[int]*Subscription),
		max_key:       0,
		mutex:         new(sync.Mutex),
		closed:        false,
	}
	return b
}

// Create a new Subscription object, which will recieve broadcasts.
func (b *Broadcaster) Subscribe() *Subscription {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.closed {
		panic("Attempt to subscribe to closed Broadcaster")
	}

	key := b.max_key
	b.max_key = key + 1
	sub := &Subscription{
		make(chan interface{}, 500),
		key,
		b,
		false,
	}
	b.subscriptions[key] = sub
	return sub
}

// Close all subscriptions (for draining), and reject new subs.
func (b *Broadcaster) Close() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.closed = true
	for _, sub := range b.subscriptions {
		sub.do_close()
	}
}

// Send data to all subscribers.
//
// Should never block.
func (b *Broadcaster) Send(data interface{}) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for _, sub := range b.subscriptions {
		open := sub.Send(data)
		if !open {
			sub.Overflowed = true
			sub.do_close()
		}
	}
}
