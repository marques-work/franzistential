package cli

import (
	"fmt"

	"github.com/marques-work/franzistential/conf"
	"github.com/marques-work/franzistential/run"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use: "franz",
		Long: " [ PIPELINE MODE ]\n\n" +
			"  In pipeline mode, `franz [flags] [--out:* DESTINATION_1] ... [--out:* DEST_N]` reads text from the\n" +
			"  pipe or STDIN and forwards it to all configured destinations.",
		RunE: func(cmd *cobra.Command, args []string) error {
			conf.Configure()
			return run.Client(opts)
		},
		Example: `
  tail -f my-server-log | franz > /dev/log \            # Forwards application server logs to
    --out:eventhub <EVENT HUB CONNECTION STRING 1> \    #   the first Event Hub endpoint,
    --out:eventhub <EVENT HUB CONNECTION STRING 2> \    #   then the second one,
    --out:stdout                                        #   and finally to STDOUT, redirecting to the system log at /dev/log


  ps -eo pid,ppid,cpu,vsize,command | franz \           # Forwards process telemetry to
    --out:eventhub <EVENT HUB CONNECTION STRING>        #   the configured Event Hub endpoint


  franz < /var/log/dmesg \                              # Forwards system kernel messages to
    --out:eventhub <EVENT HUB CONNECTION STRING>        #   the configured Event Hub endpoint`,
	}

	opts *conf.Options = conf.Instance
)

const USAGE_HEADER = `
franzistentiald [🪲]

  An event forwarding multiplexer to Kafka-esque streaming endpoints.

Synopsis:`

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	f := rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		fmt.Fprint(c.OutOrStderr(), USAGE_HEADER)
		f(c, args)
	})

	rootCmd.PersistentFlags().BoolP("help", "h", false, "Show this help message and exit")
	rootCmd.PersistentFlags().Var(newEventHubValue(opts), "out:eventhub", "Forward data to the provided Azure Event Hub `connection-uri`; multiple invocations of this flag\n  will multiplex over each connection")
	opts.Out = rootCmd.PersistentFlags().Bool("out:stdout", false, "Print data to STDOUT")
	opts.ConnectTimeout = rootCmd.PersistentFlags().Uint64P("connect-timeout-ms", "t", uint64(20*1000), "Number of `milliseconds` to allow connecting to destination before aborting message")
	opts.Silent = rootCmd.PersistentFlags().BoolP("quiet", "q", false, "``Silence logging output; does NOT affect the `--out:stdout` flag")
	opts.Trace = rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose debug output")
}
