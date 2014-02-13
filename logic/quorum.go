package model

import (
	"github.com/campadrenalin/go-deje/serial"
	"github.com/campadrenalin/go-deje/util"
)

type Quorum model.Quorum

// Serialization

func (q Quorum) ToSerial() serial.Quorum {
	return serial.Quorum{
		EventHash:  q.EventHash,
		Signatures: q.Signatures,
	}
}

func QuorumFromSerial(sq serial.Quorum) Quorum {
	return Quorum{
		EventHash:  sq.EventHash,
		Signatures: sq.Signatures,
	}
}
