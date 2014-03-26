package protocol

import (
	"encoding/json"
	"fmt"
	"github.com/campadrenalin/go-deje/logic"
	"github.com/campadrenalin/go-deje/model"
	"github.com/campadrenalin/go-deje/services"
	"log"
	"strings"
)

type Connection struct {
	Document logic.Document
	Channel  services.IRCChannel
}

func NewConnection(d logic.Document, c services.IRCChannel) Connection {
	return Connection{d, c}
}

func (p Connection) PublishEvent(ev model.Event) error {
	marshalled, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	p.Channel.Outgoing <- "deje event " + string(marshalled)
	return nil
}

func (p Connection) onEvent(evstr string) error {
	var ev model.Event
	err := json.Unmarshal([]byte(evstr), &ev)
	if err != nil {
		return err
	}
	event := logic.Event{ev, &p.Document}
	event.Register()
	return nil
}

func (p Connection) onRecv(input string) error {
	if !strings.HasPrefix(input, "deje ") {
		return fmt.Errorf(`Not a protocol message: "%s"`, input)
	}
	input = strings.TrimPrefix(input, "deje ")
	if strings.HasPrefix(input, "event ") {
		input = strings.TrimPrefix(input, "event ")
		return p.onEvent(input)
	}
	return fmt.Errorf(
		`Not a valid message type: "%s"`,
		strings.SplitN(input, " ", 2)[0],
	)
}

func (p Connection) Run(logger *log.Logger) {
	for {
		str := <-p.Channel.Incoming
		err := p.onRecv(str)
		if err != nil {
			logger.Println(err)
		}
	}
}

func (p Connection) Stop() {
}
