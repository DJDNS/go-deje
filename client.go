package deje

import (
	"github.com/DJDNS/go-deje/document"
	"github.com/jcelliott/turnpike"
)

// Contains a document and a WAMP connection.
type Client struct {
	Doc *document.Document

	onEvent  *OnEventCallback
	tpClient *turnpike.Client
}

func NewClient(topic string) Client {
	doc := document.NewDocument()
	doc.Topic = topic
	return Client{
		Doc:      &doc,
		tpClient: turnpike.NewClient(),
	}
}

// When a session connects, this callback is provided with the session
// ID value. sessionId is just a server-chosen string.
type OnConnectCallback func(sessionId string)

// Called when another peer in the same WAMP topic publishes a
// JSON-compatible event to all subscribers.
type OnEventCallback func(event interface{})

// Set callback to be executed when a successful connection has been
// made to a WAMP router.
func (c *Client) SetConnectCallback(callback OnConnectCallback) {
	c.tpClient.SetSessionOpenCallback(callback)
}

// Set callback to be executed when a published event is received.
func (c *Client) SetEventCallback(callback OnEventCallback) {
	c.onEvent = &callback
}

// Publish an event to all subscribers. Can be any JSON-compatible
// value.
func (c *Client) Publish(event interface{}) error {
	return c.tpClient.PublishExcludeMe(
		c.Doc.Topic,
		event,
	)
}

// Connect to a WAMP router. This also calls the 'connect' callback on
// success.
func (c *Client) Connect(url string) error {
	err := c.tpClient.Connect(url, "http://localhost/")
	if err != nil {
		return err
	}

	handler := func(topic string, event interface{}) {
		if c.onEvent != nil {
			(*c.onEvent)(event)
		}
	}
	return c.tpClient.Subscribe(c.Doc.Topic, handler)
}
