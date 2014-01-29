package deje

type SyncSet map[string]Sync

type Sync struct {
	EventHash    string
	Signatures   []string
	Confirmation string
}
