package state

type EventPrimitives struct {
	Id         string
	Primitives []Primitive
}

type DocumentState struct {
	Value   Container
	Applied []EventPrimitives
}

func NewDocumentState() *DocumentState {
	// We know this won't fail, so we can ignore err
	container, _ := MakeContainer(map[string]interface{}{})
	return &DocumentState{
		container,
		make([]EventPrimitives, 0),
	}
}

func (ds *DocumentState) Export() interface{} {
	return ds.Value.Export()
}
