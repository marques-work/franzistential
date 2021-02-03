package domain

import (
	"fmt"

	"github.com/Azure/azure-amqp-common-go/v3/conn"
)

type Destination interface {
	ConnectString() string
	Redacted() string
}

type EventHub struct {
	url      string
	redacted string
}

func (eh *EventHub) ConnectString() string {
	return eh.url
}

func (eh *EventHub) Redacted() string {
	return eh.redacted
}

func NewEventHub(uri string) (*EventHub, error) {
	if u, err := conn.ParsedConnectionFromStr(uri); err == nil {
		redacted := "sb://" + u.Namespace + ".servicebus.windows.net/;SharedAccessKeyName=" + u.KeyName + ";EntityPath=" + u.HubName

		return &EventHub{url: uri, redacted: redacted}, nil
	} else {
		return nil, fmt.Errorf("Malformed Event Hub connection string: %s", err.Error())
	}
}
