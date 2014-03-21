package state

import (
	"github.com/campadrenalin/go-deje/broadcast"
	"sync"
)

type PrimitiveBroadcaster struct {
	*broadcast.Broadcaster
}

func NewPrimitiveBroadcaster() PrimitiveBroadcaster {
	return PrimitiveBroadcaster{
		broadcast.NewBroadcaster(),
	}
}

func (pb PrimitiveBroadcaster) Send(p Primitive) {
	pb.Broadcaster.Send(p)
}

func (pb PrimitiveBroadcaster) Subscribe() *PrimitiveSubscription {
	ps := PrimitiveSubscription{
		pb.Broadcaster.Subscribe(),
		make(chan Primitive),
		new(sync.Mutex),
		0,
	}
	go ps.run()
	return &ps
}

type PrimitiveSubscription struct {
	sub       *broadcast.Subscription
	out       chan Primitive
	state_mut *sync.Mutex
	sending   int // Number of items being sent
}

func (ps *PrimitiveSubscription) Out() <-chan Primitive {
	return ps.out
}
func (ps *PrimitiveSubscription) Len() int {
	return ps.sub.Len() + len(ps.out) + ps.sending
}

// Secret underlying goroutine to type-assert everything to
// Primitive. We know that everything is going to be Primitive,
// as long as it's sent through the Broadcaster, but we need
// the cast to present ps.out as a <-chan Primitive.
func (ps *PrimitiveSubscription) run() {
	input := ps.sub.Out()
	for {
		ps.state_mut.Lock()
		value, ok := <-input
		ps.sending = 1
		ps.state_mut.Unlock()

		if !ok {
			close(ps.out)
			return
		}
		ps.out <- value.(Primitive)
		ps.sending = 0
	}
}

/*
func (ps *PrimitiveSubscription) Close() {
    ps.sub.Close()
}
*/
