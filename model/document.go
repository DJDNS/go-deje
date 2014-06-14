package model

// Used for serializing and deserializing docs to files.
//
// This allows us to use more complicated structures for actual
// documents, that allow for listing events or quorums by group.
type DocumentFile struct {
	Topic   string
	Events  EventSet
	Quorums QuorumSet
}

func NewDocumentFile() DocumentFile {
	return DocumentFile{
		Events:  make(EventSet),
		Quorums: make(QuorumSet),
	}
}
