package run

import (
	"context"
	"time"

	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	"github.com/marques-work/franzistential/domain"
)

// Multiplexer broadcasts a channel across Producers
type Multiplexer struct {
	egress []domain.Producer
}

// Accept ingests a single datum
func (m *Multiplexer) Accept(line string) {
	for _, pr := range m.egress {
		pr.Accept(line)
	}
}

type Producer struct {
	tx domain.Transport
}

func send(payload string, timeout uint64, hub *eventhub.Hub) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
	defer cancel()

	// send a single message into a random partition
	return hub.Send(ctx, eventhub.NewEventFromString(payload))
}

func push(payload []string, timeoutMs uint64) {}
