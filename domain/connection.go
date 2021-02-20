package domain

import (
	"fmt"
	"io"
	"os"

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
	if u, err := conn.ParsedConnectionFromStr(uri); err != nil {
		return nil, fmt.Errorf("Malformed Event Hub connection string: %s", err.Error())
	} else {
		redacted := "Endpoint=sb://" + u.Namespace + "." + u.Suffix + "/;SharedAccessKeyName=" + u.KeyName + ";SharedAccessKey=********;EntityPath=" + u.HubName
		return &EventHub{url: uri, redacted: redacted}, nil
	}
}

type IODest struct {
	stream io.WriteCloser
}

func (i *IODest) ConnectString() string {
	switch i.stream.(type) {
	case *os.File:
		return i.stream.(*os.File).Name()
	default:
		return fmt.Sprintf("%T: %v", i.stream, i.stream)
	}
}

func (i *IODest) Redacted() string {
	return i.ConnectString()
}

func NewIO(stream io.WriteCloser) (*IODest, error) {
	return &IODest{stream: stream}, nil
}
