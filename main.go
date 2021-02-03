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
	"github.com/marques-work/franzistential/domain"
	"github.com/marques-work/franzistential/logging"
	slog "gopkg.in/mcuadros/go-syslog.v2"
	format "gopkg.in/mcuadros/go-syslog.v2/format"
)

var (
	RAW     format.Format = &domain.RawFormat{}
	RFC3164 format.Format = &format.RFC3164{}
	RFC5424 format.Format = &format.RFC5424{}
	RFC6587 format.Format = &format.RFC6587{}
	DETECT  format.Format = &format.Automatic{}
)

func main() {
	conf.ConfigureAndValidate()

	if conf.ServerMode() {
		channel := make(slog.LogPartsChannel)
		handler := slog.NewChannelHandler(channel)

		server := slog.NewServer()
		server.SetFormat(RAW)

		server.SetHandler(handler)

		if "" != *conf.TcpService {
			if err := server.ListenTCP(*conf.TcpService); err != nil {
				logging.Die("Couldn't establish a TCP listener on %s: %v", *conf.TcpService, err)
			}
		}

		if "" != *conf.UdpService {
			if err := server.ListenUDP(*conf.UdpService); err != nil {
				logging.Die("Couldn't establish a UDP listener on %s: %v", *conf.UdpService, err)
			}
		}

		if "" != *conf.UnixSocket {
			if err := os.RemoveAll(*conf.UnixSocket); err != nil {
				logging.Die("Socket appears to be busy; could not remove %s: %v", *conf.UnixSocket, err)
			}

			if err := server.ListenUnixgram(*conf.UnixSocket); err != nil {
				logging.Die("Couldn't open a listening socket at %s: %v", *conf.UnixSocket, err)
			}
		}

		server.Boot()

		dest := conf.Destinations[0]
		hub, err := eventhub.NewHubFromConnectionString(dest.ConnectString())

		if err != nil {
			logging.Die("Unable to establish connection to destination %v: %v", dest, err)
		}

		// consumes parsed log output handled by server
		go func(channel slog.LogPartsChannel) {
			for logParts := range channel {
				send(logParts["msg"].(string), 20000, hub)
				if *conf.Out {
					fmt.Println(logParts["msg"].(string))
				}
			}
		}(channel)

		server.Wait()
	} else {
		input := bufio.NewScanner(os.Stdin)

		dest := conf.Destinations[0]
		hub, err := eventhub.NewHubFromConnectionString(dest.ConnectString())

		if err != nil {
			logging.Die("Unable to establish connection to destination %v: %v", dest, err)
		}

		for input.Scan() {
			line := input.Text()
			if err := send(line, *conf.SendTimeout, hub); err != nil {
				// Do we do more than this?
				logging.Warn("Failed to send event %s to Event Hub [namespace]@[queue] because: %v", input.Text(), err)
			}

			if *conf.Out {
				fmt.Println(line)
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
