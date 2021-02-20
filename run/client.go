package run

import (
	"bufio"
	"fmt"
	"log"
	"os"

	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	"github.com/marques-work/franzistential/conf"
	"github.com/marques-work/franzistential/logging"
)

// Client forwards data to the configured outputs
func Client(opts *conf.Options) error {
	logging.Debug("ROOT - opts: %v", opts)
	input := bufio.NewScanner(os.Stdin)

	if !opts.HasDestination() {
		return nil
	}

	dest := opts.Destinations[0]
	hub, err := eventhub.NewHubFromConnectionString(dest.ConnectString())

	if err != nil {
		logging.Die("Unable to establish connection to destination %v: %v", dest, err)
	}

	for input.Scan() {
		line := input.Text()
		if err := send(line, *opts.ConnectTimeout, hub); err != nil {
			// Do we do more than this?
			logging.Warn("Failed to send event %s to Event Hub [namespace]@[queue] because: %v", input.Text(), err)
		}

		if *opts.Out {
			fmt.Println(line)
		}
	}

	if err := input.Err(); err != nil {
		log.Fatalf("Unexpected error: %v\n", err)
	}
	return nil
}
