package state

import "github.com/campadrenalin/go-deje/broadcast"

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

func (pb PrimitiveBroadcaster) Subscribe() PrimitiveSubscription {
	ps := PrimitiveSubscription{
		pb.Broadcaster.Subscribe(),
		make(chan Primitive, 1),
	}
	go ps.run()
	return ps
}

type PrimitiveSubscription struct {
	sub broadcast.Subscription
	out chan Primitive
}

func (ps PrimitiveSubscription) Out() <-chan Primitive {
	return ps.out
}
func (ps PrimitiveSubscription) Len() int {
	return ps.sub.Len() + len(ps.out)
}

// Secret underlying goroutine to type-assert everything to
// Primitive. We know that everything is going to be Primitive,
// as long as it's sent through the Broadcaster, but we need
// the cast to present ps.out as a <-chan Primitive.
func (ps PrimitiveSubscription) run() {
	input := ps.sub.Out()
	value := <-input
	ps.out <- value.(Primitive)
}

/*
func (ps PrimitiveSubscription) Close() {
    ps.sub.Close()
}
*/
