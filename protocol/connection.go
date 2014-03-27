package protocol

import (
	"encoding/json"
	"fmt"
	"github.com/campadrenalin/go-deje/logic"
	"github.com/campadrenalin/go-deje/model"
	"github.com/campadrenalin/go-deje/services"
	"log"
	"strings"
	"sync"
)

type Connection struct {
	Document logic.Document
	Channel  services.IRCChannel
	closer   chan struct{}
	waiter   *sync.Mutex
}

func NewConnection(d logic.Document, c services.IRCChannel) Connection {
	return Connection{
		d,
		c,
		make(chan struct{}),
		new(sync.Mutex),
	}
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
	p.waiter.Lock()
	defer p.waiter.Unlock()
	for {
		select {
		case str := <-p.Channel.Incoming:
			err := p.onRecv(str)
			if err != nil {
				logger.Println(err)
			}
		case <-p.closer:
			logger.Println("Exiting protocol connection loop")
			return
		}
	}
}

func (p Connection) Stop() {
	close(p.closer)
	// Wait for close
	p.waiter.Lock()
	p.waiter.Unlock()
}
