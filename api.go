// This package implements the DEJE Next protocol.
// For more information, read the docs below, or README.md.
//
// The front-facing API is fairly simple, and mostly consists
// of these top-level types, and model.Document. However, the
// better you understand the underlying technology, the more
// easily you will be able to integrate it into user-facing
// software, with a minimum of behavioral surprises or conceptual
// fog.
package deje

import "github.com/campadrenalin/go-deje/model"
import "github.com/campadrenalin/go-deje/serial"

// Contains the clients for network communication and
// timestamp retrieval. Use this to create or sync to documents.
//
// You generally only want one of these per program.
type DEJEController struct {
	Timestamper TimestampService
	Networker   IRCService
}

// Get a Document based on an IRCLocation.
//
// This will create a blank document, if none exists.
// See the model.Document documentation for more information
// about how to use this object.
func (c *DEJEController) GetDocument(serial.IRCLocation) model.Document {
	return model.NewDocument()
}

// Different types of TimestampServices can be used,
// but the default implemetation makes use of the
// Bitcoin blockchain, given a bitcoind JSON-RPC location.
// Alternative timestamping services can simply be alternative
// blockchains, such as Litecoin, or something more exotic
// (if it can be made to fit the TimestampService interface).
//
// See https://en.bitcoin.it/wiki/API_reference_%28JSON-RPC%29
// for more information about this API.
type TimestampService interface {
	GetTimestampsAfter(dochash string, after model.BlockHeight)

	MakeTimestamp(dochash string, qhash string)
}

// You should generally never need a custom IRCService,
// but you can provide one, if you really want.
type IRCService interface {
	GetChannel(serial.IRCLocation) IRCChannel
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
	Location serial.IRCLocation
	Channel  chan string
}
