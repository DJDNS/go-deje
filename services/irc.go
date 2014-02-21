package services

import "github.com/campadrenalin/go-deje/model"

// You should generally never need a custom IRCService,
// but you can provide one, if you really want.
type IRCService interface {
	GetChannel(model.IRCLocation) IRCChannel
}

// An IRCChannel represents a connection to a specific
// channel on a specific IRC network. The .Channel will
// output any messages from the IRC channel, and you can
// broadcast any message by sending it to the .Channel.
//
// This provides a convenient abstraction over the underlying
// IRC implementation, allowing for transparent reuse of
// client connections.
type IRCChannel struct {
	Location model.IRCLocation
	Incoming chan string
	Outgoing chan string
}

type DummyIRCService struct{}

func (dis DummyIRCService) GetChannel(location model.IRCLocation) IRCChannel {
	return IRCChannel{
		Location: location,
		Incoming: make(chan string),
		Outgoing: make(chan string),
	}
}