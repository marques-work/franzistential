package cli

import (
	"fmt"

	"github.com/marques-work/franzistential/conf"
	"github.com/marques-work/franzistential/run"
	"github.com/spf13/cobra"
)

var (
	serverCmd = &cobra.Command{
		Use:   "daemon",
		Short: "The \"d\" in franzistentiald - runs a syslog service daemon to forward data to destinations",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Server = server
			conf.Configure()
			return run.Daemon(opts)
		},
	}

	server = &conf.ServerOptions{Parser: conf.DETECT}
)

func init() {
	serverCmd.Flags().VarP(newParserValue(server), "input-format", "F", fmt.Sprintf("Specify the `name` of the syslog input parsing format; choose one of:\n  [ %s ]; Default: guess-rfc", parserList()))
	server.UnixSocket = serverCmd.Flags().String("in:socket", "", "Listen on local unixgram (i.e., packet-based) socket `path`; socket will be created on startup")
	server.TCPService = serverCmd.Flags().String("in:tcp", "", "Listen on TCP/IP `ipaddr:portnum`, e.g. `127.0.0.1:514`; must include IP (v4 or v6) and port")
	server.UDPService = serverCmd.Flags().String("in:udp", "", "Listen on UDP/IP `ipaddr:portnum`, e.g. `127.0.0.1:514`; must include IP (v4 or v6) and port")

	rootCmd.AddCommand(serverCmd)
}