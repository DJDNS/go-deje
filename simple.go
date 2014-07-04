package deje

import (
	"errors"
	"log"

	"github.com/campadrenalin/go-deje/document"
	"github.com/campadrenalin/go-deje/util"
)

// Wraps the low-level capabilities of the basic Client to provide
// an easier, more useful API to downstream code.
type SimpleClient struct {
	client *Client
	tip    string
	logger *log.Logger
}

func NewSimpleClient(topic string, logger *log.Logger) *SimpleClient {
	raw_client := NewClient(topic)
	simple_client := &SimpleClient{&raw_client, "", logger}
	raw_client.SetEventCallback(func(event interface{}) {
		err := simple_client.onRcv(event)
		if err != nil && simple_client.logger != nil {
			simple_client.logger.Println(err)
		}
	})
	return simple_client
}

func (sc *SimpleClient) onRcv(event interface{}) error {
	map_ev, ok := event.(map[string]interface{})
	if !ok {
		return errors.New("Non-{} message")
	}
	evtype, ok := map_ev["type"].(string)
	if !ok {
		return errors.New("Message with no 'type' param")
	}

	doc := sc.GetDoc()
	switch evtype {
	case "01-request-tip":
		sc.PublishTip()
	case "01-publish-tip":
		hash, ok := map_ev["tip_hash"].(string)
		if !ok {
			return errors.New("Message with bad 'tip_hash' param")
		}
		if sc.tip != hash {
			sc.RequestHistory()
		}
	case "01-request-history":
		sc.PublishHistory()
	case "01-publish-history":
		// This is intentionally structured so that the
		// coverage tests will be helpful for catching all
		// possible circumstances.
		history, ok := map_ev["history"].([]interface{})
		if !ok {
			return errors.New("History message with bad 'history' param")
		}
		for _, serial_event := range history {
			doc_ev := doc.NewEvent("")
			err := util.CloneMarshal(serial_event, &doc_ev)
			if err != nil {
				return err
			}
			doc_ev.Register()
		}
		hash, ok := map_ev["tip_hash"].(string)
		if !ok {
			return errors.New("Message with bad 'tip_hash' param")
		}
		tip_event, ok := doc.Events[hash]
		if !ok {
			return errors.New("Unknown event " + hash)
		}
		err := tip_event.Goto()
		if err != nil {
			return err
		}
		sc.tip = hash
		sc.logger.Printf("Updated content: %#v", sc.Export())
	default:
		return errors.New("Unfamiliar message type: '" + evtype + "'")
	}
	return nil
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

// Navigate the Document to an Event, and promote it as the tip.
func (sc *SimpleClient) Promote(ev document.Event) error {
	if err := ev.Goto(); err != nil {
		return err
	}
	sc.tip = ev.Hash()
	return sc.PublishTip()
}

// Get the Document object owned by this Client.
func (sc *SimpleClient) GetDoc() *document.Document {
	return sc.client.Doc
}

// Return the current contents of the document.
func (sc *SimpleClient) Export() interface{} {
	return sc.client.Doc.State.Export()
}
