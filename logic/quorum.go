package logic

import (
	"github.com/campadrenalin/go-deje/model"
)

type Quorum struct {
	model.Quorum
	Doc *Document
}

func (doc *Document) NewQuorum(evhash string) Quorum {
	return Quorum{
		model.NewQuorum(evhash),
		doc,
	}
}