package deje

import "github.com/campadrenalin/go-deje/document"

// Wraps the low-level capabilities of the basic Client to provide
// an easier, more useful API to downstream code.
type SimpleClient struct {
	client Client
}

func NewSimpleClient(topic string) SimpleClient {
	return SimpleClient{NewClient(topic)}
}

// Connect and immediately request the tip event hash.
func (sc *SimpleClient) Connect(url string) error {
	return sc.client.Connect(url)
}

// Get the Document object owned by this Client.
func (sc *SimpleClient) GetDoc() *document.Document {
	return sc.client.Doc
}

// Return the current contents of the document.
func (sc *SimpleClient) Export() interface{} {
	return sc.client.Doc.State.Export()
}
