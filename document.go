package deje

type Document struct {
	Channel   IRCLocation
	Downloads map[string]string
	Events    EventSet
	Syncs     map[string]Sync
}
