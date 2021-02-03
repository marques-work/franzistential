package run

import (
	"fmt"
	"os"

	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	"github.com/marques-work/franzistential/conf"
	"github.com/marques-work/franzistential/logging"
	slog "gopkg.in/mcuadros/go-syslog.v2"
)

// Daemon starts a syslog server
func Daemon(opts *conf.Options) error {
	logging.Debug("SERVER - opts: %v", opts)

	config := opts.Server
	channel := make(slog.LogPartsChannel)
	handler := slog.NewChannelHandler(channel)

	server := slog.NewServer()
	server.SetFormat(config.Parser)

	server.SetHandler(handler)

	if "" != *config.TCPService {
		if err := server.ListenTCP(*config.TCPService); err != nil {
			logging.Die("Couldn't establish a TCP listener on %s: %v", *config.TCPService, err)
		}
	}

	if "" != *config.UDPService {
		if err := server.ListenUDP(*config.UDPService); err != nil {
			logging.Die("Couldn't establish a UDP listener on %s: %v", *config.UDPService, err)
		}
	}

	if "" != *config.UnixSocket {
		if err := os.RemoveAll(*config.UnixSocket); err != nil {
			logging.Die("Socket appears to be busy; could not remove %s: %v", *config.UnixSocket, err)
		}

		if err := server.ListenUnixgram(*config.UnixSocket); err != nil {
			logging.Die("Couldn't open a listening socket at %s: %v", *config.UnixSocket, err)
		}
	}

	server.Boot()

	if !opts.HasDestination() {
		return nil
	}

	dest := opts.Destinations[0]
	hub, err := eventhub.NewHubFromConnectionString(dest.ConnectString())

	if err != nil {
		logging.Die("Unable to establish connection to destination %v: %v", dest, err)
	}

	// consumes parsed log output handled by server
	go func(channel slog.LogPartsChannel) {
		for logParts := range channel {
			send(logParts["msg"].(string), 20000, hub)
			if *opts.Out {
				fmt.Println(logParts["msg"].(string))
			}
		}
	}(channel)

	server.Wait()

	return nil
}
