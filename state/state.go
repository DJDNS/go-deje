package state

type DocumentState struct {
	Value Container
	bcast PrimitiveBroadcaster
}

func NewDocumentState() *DocumentState {
	// We know this won't fail, so we can ignore err
	container, _ := MakeContainer(map[string]interface{}{})
	return &DocumentState{
		container,
		NewPrimitiveBroadcaster(),
	}
}

func (ds *DocumentState) Subscribe() PrimitiveSubscription {
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
