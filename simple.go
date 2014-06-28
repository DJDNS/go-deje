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
		evtype, ok := map_ev["type"].(string)
		if !ok {
			return
		}

		switch evtype {
		case "01-request-tip":
			simple_client.PublishTip()
		case "01-request-history":
			simple_client.PublishHistory()
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

func (sc *SimpleClient) PublishHistory() error {
	response := map[string]interface{}{
		"type":     "01-publish-history",
		"tip_hash": sc.tip,
	}
	doc := sc.GetDoc()
	ev, ok := doc.Events[sc.tip]
	if !ok {
		response["error"] = "not-found"
		return sc.client.Publish(response)
	}

	history, ok := ev.GetHistory()
	if !ok {
		response["error"] = "root-not-found"
		return sc.client.Publish(response)
	}
	response["history"] = history
	return sc.client.Publish(response)
}

// Get the Document object owned by this Client.
func (sc *SimpleClient) GetDoc() *document.Document {
	return sc.client.Doc
}

// Return the current contents of the document.
func (sc *SimpleClient) Export() interface{} {
	return sc.client.Doc.State.Export()
}
