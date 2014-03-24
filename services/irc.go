package services

import "github.com/campadrenalin/go-deje/model"

var CHANNEL_BUFFER_SIZE = 5

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

func NewIRCChannel(location model.IRCLocation) IRCChannel {
	return IRCChannel{
		Location: location,
		Incoming: make(chan string, CHANNEL_BUFFER_SIZE),
		Outgoing: make(chan string, CHANNEL_BUFFER_SIZE),
	}
}

type DummyIRCService struct{}

func (dis DummyIRCService) GetChannel(location model.IRCLocation) IRCChannel {
	return NewIRCChannel(location)
}
