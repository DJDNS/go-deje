package deje

type Event struct {
	ParentHash  SHA1Hash               `json:"parent"`
	HandlerName string                 `json:"handler"`
	Arguments   map[string]interface{} `json:"args"`
}

type EventSet map[SHA1Hash]Event

func NewEvent(hname string) Event {
	return Event{
		ParentHash:  SHA1Hash{},
		HandlerName: hname,
		Arguments:   make(map[string]interface{}),
	}
}

func (e *Event) SetParent(p Event) error {
	hash, err := HashObject(p)
	if err != nil {
		return err
	}
	e.ParentHash = hash
	return nil
}

func (s EventSet) Register(event Event) {
	hash, _ := HashObject(event)
	s[hash] = event
}

func (s EventSet) GetRoot(tip Event) (event Event, ok bool) {
	event = tip
	ok = true
	for (event.ParentHash != SHA1Hash{}) {
		event, ok = s[event.ParentHash]
		if !ok {
			return
		}
	}
	return
}
