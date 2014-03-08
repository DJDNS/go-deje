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
	container, _ := MakeContainer(map[string]interface{}{})
	ds.Value = container
}

func (ds *DocumentState) Export() interface{} {
	return ds.Value.Export()
}
