package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"time"

	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	"github.com/marques-work/franzistential/conf"
	"github.com/marques-work/franzistential/domain"
	"github.com/marques-work/franzistential/logging"
	sls "gopkg.in/mcuadros/go-syslog.v2"
	format "gopkg.in/mcuadros/go-syslog.v2/format"
)

func main() {
	conf.ConfigureAndValidate()

	if conf.ServerMode() {
		server := sls.NewServer()

		if *conf.RawForward {
			server.SetFormat(&domain.RawFormat{})
		}

		ch := make(chan format.LogParts)
		server.SetHandler(sls.NewChannelHandler(ch))
	} else {
		input := bufio.NewScanner(os.Stdin)

		dest := conf.Destinations[0]
		hub, err := eventhub.NewHubFromConnectionString(dest.ConnectString())

		if err != nil {
			logging.Die("Unable to establish connection: %v", err)
			return
		}

		for input.Scan() {
			if err := send(input.Text(), *conf.SendTimeout, hub); err != nil {
				// Do we do more than this?
				logging.Warn("Failed to send event %s to Event Hub [namespace]@[queue] because: %v", input.Text(), err)
			}
		}

		if err := input.Err(); err != nil {
			log.Fatalf("Unexpected error: %v\n", err)
		}
	}
}

func send(payload string, timeout uint64, hub *eventhub.Hub) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
	defer cancel()

	// send a single message into a random partition
	return hub.Send(ctx, eventhub.NewEventFromString(payload))
}
