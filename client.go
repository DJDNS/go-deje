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
