package domain

import (
	"context"
	"io"
	"time"

	eventhub "github.com/Azure/azure-event-hubs-go/v3"
)

type Transport interface {
	Send(payload []string, timeoutMs uint64) error
	Endpoint() string

	io.Closer
}

type EventHubTransport struct {
	hub *eventhub.Hub
	ep  string
}

func (e *EventHubTransport) Send(payload []string, timeout uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
	defer cancel()

	return e.hub.SendBatch(ctx, e.wrap(payload))
}

func (e *EventHubTransport) Endpoint() string {
	return e.ep
}

func (e *EventHubTransport) Close() error {
	return nil
}

func (e *EventHubTransport) wrap(payload []string) *eventhub.EventBatchIterator {
	size := len(payload)
	events := make([]*eventhub.Event, size, size)

	for i, line := range payload {
		events[i] = eventhub.NewEventFromString(line)
	}

	return eventhub.NewEventBatchIterator(events...)
}

func NewEventHubTx(dest Destination) (*EventHubTransport, error) {
	hub, err := eventhub.NewHubFromConnectionString(dest.ConnectString())

	if err != nil {
		return nil, err
	}

	return &EventHubTransport{hub: hub, ep: dest.Redacted()}, nil
}
