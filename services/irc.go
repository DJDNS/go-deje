package services

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/campadrenalin/go-deje/model"
	external "github.com/thoj/go-ircevent"
)

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

func randomNick() string {
	buffer := make([]byte, 20)
	rand.Read(buffer)
	encoder := base64.StdEncoding
	enc_out := make([]byte, encoder.EncodedLen(len(buffer)))
	encoder.Encode(enc_out, buffer)
	return string(enc_out)
}

type DummyIRCService struct{}

func (dis DummyIRCService) GetChannel(location model.IRCLocation) IRCChannel {
	return NewIRCChannel(location)
}

type RealIRCService struct {
}

func (dis RealIRCService) GetChannel(location model.IRCLocation) IRCChannel {
	channel := NewIRCChannel(location)

	nick := randomNick()
	backend := external.IRC(nick, nick)
	err := backend.Connect(location.GetHostPort())
	if err != nil {
		panic(err)
	}
	backend.AddCallback("JOIN", func(e *external.Event) {
		msg := e.Raw
		channel.Incoming <- msg
	})
	backend.Join("#" + location.Channel)

	// Incoming
	backend.AddCallback("PRIVMSG", func(e *external.Event) {
		msg := e.Message()
		channel.Incoming <- msg
	})

	// Outgoing
	go func() {
		for {
			line := <-channel.Outgoing
			backend.Privmsg("#"+location.Channel, line)
		}
	}()
	return channel
}
