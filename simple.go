package deje

// Wraps the low-level capabilities of the basic Client to provide
// an easier, more useful API to downstream code.
type SimpleClient struct {
	client Client
}

func NewSimpleClient(topic string) SimpleClient {
	return SimpleClient{NewClient(topic)}
}
