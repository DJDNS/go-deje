package deje

import (
	"github.com/campadrenalin/go-deje/logic"
	"github.com/jcelliott/turnpike"
)

// Contains a document and one or more WAMP connections.
type Client struct {
	Doc *logic.Document

	tpClient *turnpike.Client
}

func NewClient(topic string) Client {
	doc := logic.NewDocument()
	doc.Topic = topic
	return Client{
		Doc:      &doc,
		tpClient: turnpike.NewClient(),
	}
}

// When a session connects, this client is provided with the session
// ID value. This is just a server-chosen string.
type OnConnectCallback func(sessionId string)

// Set callback to be executed when a successful connection has been
// made to a WAMP router.
func (c *Client) SetConnectCallback(callback OnConnectCallback) {
	c.tpClient.SetSessionOpenCallback(callback)
}

// Connect to a WAMP router. This also calls the 'connect' callback on
// success.
func (c *Client) Connect(url string) error {
	return c.tpClient.Connect(url, "http://localhost/")
}
