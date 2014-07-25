package state

// Represents the state of a document. It can be modified by applying
// Primitives to it, which make simple replacements or deletions to
// the contents at specified locations.
//
// These transformations not only alter the DocumentState's .Value
// field, they also call the OnPrimitive callback, if it is set for
// this DocumentState.
type DocumentState struct {
	Value       Container
	onPrimitive OnPrimitiveCallback
}

func NewDocumentState() *DocumentState {
	// We know this won't fail, so we can ignore err
	container, _ := makeContainer(map[string]interface{}{})
	return &DocumentState{container, nil}
}

// Construct and apply a Primitive that completely resets the Value
// of the DocumentState to an empty JSON {}.
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
	ds.onPrimitive = c
}

// Apply a Primitive such that the callback (if set) is run.
//
// Always preferable to p.Apply(ds), which does not run the callback.
func (ds *DocumentState) Apply(p Primitive) error {
	err := p.Apply(ds)
	if err != nil {
		return err
	}
	if ds.onPrimitive != nil {
		ds.onPrimitive(p)
	}
	return nil
}

// Return the raw, JSON-ic value of the DocumentState.
func (ds *DocumentState) Export() interface{} {
	return ds.Value.Export()
}
