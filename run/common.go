package run

import (
	"context"
	"time"

	eventhub "github.com/Azure/azure-event-hubs-go/v3"
)

func send(payload string, timeout uint64, hub *eventhub.Hub) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
	defer cancel()

	// send a single message into a random partition
	return hub.Send(ctx, eventhub.NewEventFromString(payload))
}
