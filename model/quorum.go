package model

// Represents a complete set of approvals for an event.
// Quorums act as bridges between events and timestamps,
// indicating that an event was both common knowledge and
// considered a valid event chain (among others) at one
// time (the timestamp provides the time information).
type Quorum struct {
	EventHash  string            `json:"event_hash"`
	Signatures map[string]string `json:"sigs"`
}

type QuorumSet map[string]Quorum
