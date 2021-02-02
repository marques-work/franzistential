package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	"github.com/marques-work/franzistential/conf"
	sls "gopkg.in/mcuadros/go-syslog.v2"
)

func main() {
	var sendTimeout uint64 = 20 * 1000 // move to configuration
	connStr := "Endpoint=sb://namespace.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=superSecret1234=;EntityPath=hubName"

	conf.ConfigureAndValidate()

	sls.NewServer()
	input := bufio.NewScanner(os.Stdin)

	hub, err := eventhub.NewHubFromConnectionString(connStr)

	if err != nil {
		fmt.Println(err)
		return
	}

	for input.Scan() {
		if err := send(input.Text(), sendTimeout, hub); err != nil {
			// Doo we do more than this?
			fmt.Fprintf(os.Stderr, "Failed to send event %s to Event Hub [namespace]@[queue] because: %s", input.Text(), err.Error())
		}
	}

	if err := input.Err(); err != nil {
		log.Fatal(err)
	}
}

func send(payload string, timeout uint64, hub *eventhub.Hub) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
	defer cancel()

	// send a single message into a random partition
	return hub.Send(ctx, eventhub.NewEventFromString(payload))
}
