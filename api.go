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

import (
	"github.com/campadrenalin/go-deje/logic"
	"github.com/campadrenalin/go-deje/model"
	"github.com/campadrenalin/go-deje/services"
)

// Contains the clients for network communication and
// timestamp retrieval. Use this to create or sync to documents.
//
// You generally only want one of these per program.
type DEJEController struct {
	Timestamper services.TimestampService
	Networker   services.IRCService
}

func NewDEJEController() *DEJEController {
	return &DEJEController{
		Timestamper: services.DummyTimestampService{},
		Networker:   services.DummyIRCService{},
	}
}

// Get a Document based on an IRCLocation.
//
// This will create a blank document, if none exists.
// See the model.Document documentation for more information
// about how to use this object.
func (c *DEJEController) GetDocument(loc model.IRCLocation) logic.Document {
	doc := logic.NewDocument()
	doc.Channel = loc

	return doc
}
