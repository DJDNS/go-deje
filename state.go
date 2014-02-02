package deje

import "errors"

type DocumentState struct {
	Version string
	Content JSONObject
}

func NewDocumentState() DocumentState {
	return DocumentState{
		Version: "",
		Content: make(JSONObject),
	}
}

func (ds DocumentState) GetProperty(name string, s interface{}) error {
	data, ok := ds.Content[name]
	if !ok {
		return errors.New("Document does not have " + name + " property")
	}

	return CloneMarshal(data, s)
}