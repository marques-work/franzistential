package domain

type Destination interface {
	ConnectString() string
}

type EventHub struct {
	url string
}

func (eh *EventHub) ConnectString() string {
	return eh.url
}

func NewEventHub(url string) *EventHub {
	return &EventHub{url: url}
}
