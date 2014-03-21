package state

import "github.com/campadrenalin/go-deje/broadcast"

type DocumentState struct {
	Value Container
	bcast *broadcast.Broadcaster
}

func NewDocumentState() *DocumentState {
	// We know this won't fail, so we can ignore err
	container, _ := MakeContainer(map[string]interface{}{})
	return &DocumentState{
		container,
		broadcast.NewBroadcaster(),
	}
}

func (ds *DocumentState) Subscribe() *broadcast.Subscription {
	return ds.bcast.Subscribe()
}

func (ds *DocumentState) Reset() {
	// We know this won't fail, so we can ignore err
	p := &SetPrimitive{
		[]interface{}{},
		map[string]interface{}{},
	}
	ds.Apply(p)
}

// Apply a Primitive such that it is broadcast to
// all subscribers. Always preferable to p.Apply(ds),
// which does not broadcast.
func (ds *DocumentState) Apply(p Primitive) error {
	err := p.Apply(ds)
	if err != nil {
		return err
	}
	ds.bcast.Send(p)
	return nil
}

func (ds *DocumentState) Export() interface{} {
	return ds.Value.Export()
}
