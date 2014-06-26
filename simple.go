package deje

import "github.com/campadrenalin/go-deje/document"

// Wraps the low-level capabilities of the basic Client to provide
// an easier, more useful API to downstream code.
type SimpleClient struct {
	client *Client
	tip    string
}

func NewSimpleClient(topic string) *SimpleClient {
	raw_client := NewClient(topic)
	simple_client := &SimpleClient{&raw_client, ""}
	raw_client.SetEventCallback(func(event interface{}) {
		map_ev, ok := event.(map[string]interface{})
		if !ok {
			return
		}
		if map_ev["type"] == "01-request-tip" {
			simple_client.PublishTip()
		}
	})
	return simple_client
}

// Connect and immediately request the tip event hash.
func (sc *SimpleClient) Connect(url string) error {
	err := sc.client.Connect(url)
	if err != nil {
		return err
	}
	return sc.RequestTip()
}

func (sc *SimpleClient) RequestTip() error {
	return sc.client.Publish(map[string]interface{}{
		"type": "01-request-tip",
	})
}

func (sc *SimpleClient) PublishTip() error {
	return sc.client.Publish(map[string]interface{}{
		"type":     "01-publish-tip",
		"tip_hash": sc.tip,
	})
}

func (sc *SimpleClient) RequestHistory() error {
	return sc.client.Publish(map[string]interface{}{
		"type": "01-request-history",
	})
}

// Get the Document object owned by this Client.
func (sc *SimpleClient) GetDoc() *document.Document {
	return sc.client.Doc
}

// Return the current contents of the document.
func (sc *SimpleClient) Export() interface{} {
	return sc.client.Doc.State.Export()
}
