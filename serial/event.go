package serial

type Event struct {
	ParentHash  string                 `json:"parent"`
	HandlerName string                 `json:"handler"`
	Arguments   map[string]interface{} `json:"args"`
}

type EventSet map[string]Event

func NewEvent(hname string) Event {
	return Event{
		ParentHash:  "",
		HandlerName: hname,
		Arguments:   make(map[string]interface{}),
	}
}
