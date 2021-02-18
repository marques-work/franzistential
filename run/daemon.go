package run

import (
	"fmt"
	"os"

	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	"github.com/marques-work/franzistential/conf"
	"github.com/marques-work/franzistential/domain"
	"github.com/marques-work/franzistential/logging"
	slog "gopkg.in/mcuadros/go-syslog.v2"
	"gopkg.in/mcuadros/go-syslog.v2/format"
)

var (
	RAW     format.Format = &domain.RawFormat{}
	RFC3164 format.Format = &format.RFC3164{}
	RFC5424 format.Format = &format.RFC5424{}
	RFC6587 format.Format = &format.RFC6587{}
	DETECT  format.Format = &format.Automatic{}
)

func Daemon(opts *conf.Options) error {
	logging.Debug("SERVER - opts: %v", opts)

	so := opts.Server
	channel := make(slog.LogPartsChannel)
	handler := slog.NewChannelHandler(channel)

	server := slog.NewServer()
	server.SetFormat(RAW)

	server.SetHandler(handler)

	if "" != *so.TCPService {
		if err := server.ListenTCP(*so.TCPService); err != nil {
			logging.Die("Couldn't establish a TCP listener on %s: %v", *so.TCPService, err)
		}
	}

	if "" != *so.UDPService {
		if err := server.ListenUDP(*so.UDPService); err != nil {
			logging.Die("Couldn't establish a UDP listener on %s: %v", *so.UDPService, err)
		}
	}

	if "" != *so.UnixSocket {
		if err := os.RemoveAll(*so.UnixSocket); err != nil {
			logging.Die("Socket appears to be busy; could not remove %s: %v", *so.UnixSocket, err)
		}

		if err := server.ListenUnixgram(*so.UnixSocket); err != nil {
			logging.Die("Couldn't open a listening socket at %s: %v", *so.UnixSocket, err)
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
