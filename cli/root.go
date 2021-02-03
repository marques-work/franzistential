package cli

import (
	"github.com/marques-work/franzistential/conf"
	"github.com/marques-work/franzistential/run"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "franz",
		Short: "franzistentiald forwards messages to Kafka-esque streaming destinations",
		RunE: func(cmd *cobra.Command, args []string) error {
			conf.Configure()
			return run.Client(opts)
		},
	}

	opts *conf.Options = conf.Instance
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().Var(newEventHubValue(opts), "out:eventhub", "Forward data to the provided Azure Event Hub `connection-uri`; multiple invocations of this flag will multiplex over each connection")
	opts.Out = rootCmd.PersistentFlags().Bool("out:stdout", false, "Print data to STDOUT")
	opts.SendTimeout = rootCmd.PersistentFlags().Uint64P("sendTimeoutMs", "t", uint64(20*1000), "Number of `milliseconds` before aborting message to destination (default: 20000)")
	opts.Silent = rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Silence logging output; ``does NOT affect the `-out:stdout` flag")
	opts.Trace = rootCmd.PersistentFlags().BoolP("debug", "v", false, "Verbose debug output")
}
