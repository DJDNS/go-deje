package deje

import (
	"errors"
	"log"
	"sort"

	"github.com/DJDNS/go-deje/document"
	"github.com/DJDNS/go-deje/state"
	"github.com/DJDNS/go-deje/util"
)

// Wraps the low-level capabilities of the basic Client to provide
// an easier, more useful API to downstream code.
type SimpleClient struct {
	client *Client
	tip    string
	logger *log.Logger
}

// Unless you want to manually specify router URL and topic separately,
// you should probably use Open() instead of NewSimpleClient().
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

// The preferred way to create SimpleClients. Handles the Connect() call, and
// uses GetRouterAndTopic() to turn a single deje://... URL into a router URL
// and topic.
func Open(deje_url string, logger *log.Logger, cb state.OnPrimitiveCallback) (*SimpleClient, error) {
	router, topic, err := GetRouterAndTopic(deje_url)
	if err != nil {
		return nil, err
	}

	sc := NewSimpleClient(topic, logger)
	sc.SetPrimitiveCallback(cb)
	if err = sc.Connect(router); err != nil {
		return nil, err
	}
	return sc, nil
}

func (sc *SimpleClient) rcvEventList(parent map[string]interface{}, key string) error {
	doc := sc.GetDoc()
	events, ok := parent[key].([]interface{})
	if !ok {
		return errors.New("Message with bad '" + key + "' param")
	}
	for _, serial_event := range events {
		doc_ev := doc.NewEvent("")
		err := util.CloneMarshal(serial_event, &doc_ev)
		if err != nil {
			return err
		}
		doc_ev.Register()
	}
	return nil
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
	case "01-request-events":
		sc.PublishEvents()
	case "01-publish-events":
		return sc.rcvEventList(map_ev, "events")
	case "01-request-timestamps":
		sc.PublishTimestamps()
	case "01-publish-timestamps":
		ts, ok := map_ev["timestamps"].([]interface{})
		if !ok {
			return errors.New("Message with bad 'timestamps' param")
		}
		ts_strings := make([]string, len(ts))
		for i, timestamp := range ts {
			ts_string, ok := timestamp.(string)
			if !ok {
				return errors.New("Message with bad 'timestamps' param")
			}
			ts_strings[i] = ts_string
		}
		doc.Timestamps = ts_strings

		var unfamiliar bool
		for _, ts_string := range doc.Timestamps {
			if _, ok := doc.Events[ts_string]; !ok {
				unfamiliar = true
			}
		}
		if unfamiliar {
			sc.RequestEvents()
		}
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
	return sc.RequestTimestamps()
}

func (sc *SimpleClient) RequestEvents() error {
	return sc.client.Publish(map[string]interface{}{
		"type": "01-request-events",
	})
}

func (sc *SimpleClient) PublishEvents() error {
	doc := sc.GetDoc()
	hashes := make([]string, len(doc.Events))
	events := make([]*document.Event, len(doc.Events))

	// Provide events in hash-sorted order
	var i int
	for hash := range doc.Events {
		hashes[i] = hash
		i++
	}
	sort.Strings(hashes)
	for i, hash := range hashes {
		events[i] = doc.Events[hash]
	}

	return sc.client.Publish(map[string]interface{}{
		"type":   "01-publish-events",
		"events": events,
	})
}

func (sc *SimpleClient) RequestTimestamps() error {
	return sc.client.Publish(map[string]interface{}{
		"type": "01-request-timestamps",
	})
}

func (sc *SimpleClient) PublishTimestamps() error {
	return sc.client.Publish(map[string]interface{}{
		"type":       "01-publish-timestamps",
		"timestamps": sc.GetDoc().Timestamps,
	})
}

// Navigate the Document to an Event, and promote it as the tip.
func (sc *SimpleClient) Promote(ev document.Event) error {
	doc := sc.GetDoc()

	if err := ev.Goto(); err != nil {
		return err
	}
	sc.tip = ev.Hash()

	doc.Timestamps = append(doc.Timestamps, ev.Hash())
	return sc.PublishTimestamps()
}

// Set a callback for when primitives are applied to the document state.
func (sc *SimpleClient) SetPrimitiveCallback(c state.OnPrimitiveCallback) {
	sc.GetDoc().State.SetPrimitiveCallback(c)
}

// Get the Document object owned by this Client.
func (sc *SimpleClient) GetDoc() *document.Document {
	return sc.client.Doc
}

// Return the current contents of the document.
func (sc *SimpleClient) Export() interface{} {
	return sc.client.Doc.State.Export()
}
