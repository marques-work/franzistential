package cli

import (
	"fmt"

	"github.com/marques-work/franzistential/conf"
	"github.com/marques-work/franzistential/run"
	"github.com/spf13/cobra"
)

var (
	serverCmd = &cobra.Command{
		Use: "daemon",
		Long: " [ DAEMON MODE ]\n\n" +
			"  The \"d\" in franzistentiald.\n\n" +
			"  In daemon mode, `franz daemon [flags] [--out:* DESTINATION_1] ... [--out:* DEST_N]` runs as a syslog\n" +
			"  server daemon and forwards all received data to the configured destinations.",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Server = server
			conf.Configure()
			return run.Daemon(opts)
		},
		Example: `
  franz daemon --quiet \                                # Run as a syslog daemon (silencing trace, warn, and error messages)
    --in:tcp 0.0.0.0:5140 \                             #   listening on all interfaces bound to TCP port 5140
    --in:socket /var/run/franz.sock \                   #   and on a UNIXGRAM domain socket at /var/run/franz.sock
    --input-format rfc3164 \                            #   parsing the incoming data against the RFC3164 BSD format
    --out:eventhub <EVENT HUB CONNECTION STRING 1>      #   and forwarding everything to the configured Event Hub endpoint


  franz daemon \                                        # Run as a syslog daemon
    --in:udp 127.0.0.1:1514 \                           #   listening on localhost bound to UDP port 5140
    --input-format raw \                                #   allowing data to pass through as-is (i.e., unaltered)
    --out:eventhub <EVENT HUB CONNECTION STRING 1>      #   and forwarding everything to the configured Event Hub endpoint`,
	}

	server = &conf.ServerOptions{Parser: conf.DETECT}
)

func init() {
	serverCmd.Flags().VarP(newParserValue(server), "input-format", "F", fmt.Sprintf("Specify the `name` of the syslog input parsing format; choose one of:\n  [ %s ]\n\n  (default guess-rfc)", parserList()))
	server.UnixSocket = serverCmd.Flags().String("in:socket", "", "Listen on a local unixgram (packet-based) socket `path`; socket will be created on startup")
	server.TCPService = serverCmd.Flags().String("in:tcp", "", "Listen on TCP/IP `ipaddr:portnum`, e.g. `127.0.0.1:514`; must include IP (v4 or v6) and port")
	server.UDPService = serverCmd.Flags().String("in:udp", "", "Listen on UDP/IP `ipaddr:portnum`, e.g. `127.0.0.1:514`; must include IP (v4 or v6) and port")

	rootCmd.AddCommand(serverCmd)
}
