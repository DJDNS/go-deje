package protocol

import (
	"encoding/json"
	"github.com/campadrenalin/go-deje/logic"
	"github.com/campadrenalin/go-deje/model"
	"github.com/campadrenalin/go-deje/services"
)

type Connection struct {
	Channel services.IRCChannel
}

func NewConnection(d logic.Document, c services.IRCChannel) Connection {
	return Connection{c}
}

func (p Connection) PublishEvent(ev model.Event) error {
	marshalled, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	p.Channel.Outgoing <- "deje event " + string(marshalled)
	return nil
}
