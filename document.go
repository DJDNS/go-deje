package deje

type Document struct {
	Events EventSet
	Syncs  map[string]Sync
}
