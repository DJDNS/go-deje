package logic

import (
	"github.com/campadrenalin/go-deje/model"
)

type Quorum struct {
	model.Quorum
	Doc *Document
}
type QuorumSet map[string]*Quorum

func (doc *Document) NewQuorum(evhash string) Quorum {
	return Quorum{
		model.NewQuorum(evhash),
		doc,
	}
}

// Register with the Doc. This stores it in a hash-based location,
// so do not make changes to an Event after it has been registered.
func (q *Quorum) Register() {
	key := q.GetKey()
	q.Doc.Quorums[key] = q

	group_key := q.GetGroupKey()
	group, ok := q.Doc.QuorumsByEvent[group_key]
	if !ok {
		group = make(QuorumSet)
		q.Doc.QuorumsByEvent[group_key] = group
	}
	group[key] = q
}
