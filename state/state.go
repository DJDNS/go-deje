package state

type DocumentState struct {
	Value Container
}

func NewDocumentState() *DocumentState {
	// We know this won't fail, so we can ignore err
	container, _ := MakeContainer(map[string]interface{}{})
	return &DocumentState{
		container,
	}
}

func (ds *DocumentState) Reset() {
	// We know this won't fail, so we can ignore err
	p := &SetPrimitive{
		[]interface{}{},
		map[string]interface{}{},
	}
	ds.Apply(p)
}

// An optional callback to be called for every primitive applied to
// a DocumentState object. Will always be called in the same order,
// in the same goroutine, as the Primitive application itself.
type OnPrimitiveCallback func(primitive Primitive)

// Set the OnPrimitiveCallback for this DocumentState.
func (ds *DocumentState) SetPrimitiveCallback(c OnPrimitiveCallback) {
}

// Apply a Primitive such that it is broadcast to
// all subscribers. Always preferable to p.Apply(ds),
// which does not broadcast.
func (ds *DocumentState) Apply(p Primitive) error {
	err := p.Apply(ds)
	if err != nil {
		return err
	}
	//ds.bcast.Send(p)
	return nil
}

func (ds *DocumentState) Export() interface{} {
	return ds.Value.Export()
}
