package conf

import (
	"fmt"

	"github.com/marques-work/franzistential/domain"
	"github.com/marques-work/franzistential/logging"
	"gopkg.in/mcuadros/go-syslog.v2/format"
)

// Instance is the default Options value
var Instance = &Options{}

// Configure initializes and validates the default Options
func Configure() {
	Instance.BindAndValidate(&consoleOpts{
		Trace:  &logging.Trace,
		Silent: &logging.Silent,
	})
}

// Options represents configuration
type Options struct {
	Out          *bool
	Destinations []domain.Destination
	SendTimeout  *uint64

	// logging
	Silent *bool
	Trace  *bool

	Server *ServerOptions
}

// ServerMode returns whether or not to start a daemon
func (o *Options) ServerMode() bool {
	return nil != o.Server
}

// HasDestination returns whether or not there is at least one output destination configured
func (o *Options) HasDestination() bool {
	return *o.Out || len(o.Destinations) > 0
}

// BindAndValidate binds config and validates
func (o *Options) BindAndValidate(console *consoleOpts) {
	console.Bind(o)

	if console.Invalid() {
		logging.Die("Both `--debug` and `--quiet` cannot be simultaneously set; pick one")
	}

	if !o.HasDestination() {
		logging.Warn("You have not set any destinations and are not printing to STDOUT; this will act like a sink, which might be a configuration error.")
	}

	if o.ServerMode() {
		tcp := *o.Server.TCPService

		if "" != tcp {
			validateNetworkListenerConfig("tcp", tcp)
		}

		udp := *o.Server.UDPService

		if "" != udp {
			validateNetworkListenerConfig("udp", udp)
		}
	}
}

func (o *Options) String() string {
	return "Options {\n" +
		fmt.Sprintf("  Out: %t\n", *o.Out) +
		fmt.Sprintf("  Destinations (%d)%v\n", len(o.Destinations), o.dests()) +
		fmt.Sprintf("  SendTimeout: %d ms\n", *o.SendTimeout) +
		fmt.Sprintf("  Server: %v\n", o.Server) +
		fmt.Sprintf("  Trace: %t\n", *o.Trace) +
		fmt.Sprintf("  Silent: %t\n", *o.Silent) +
		"}\n"
}

func (o *Options) dests() string {
	if len(o.Destinations) == 0 {
		return ""
	}

	s := ":"
	for _, d := range o.Destinations {
		s += fmt.Sprintf("\n    - %s", d.Redacted())
	}
	return s
}

// ServerOptions represents daemon configuration
type ServerOptions struct {
	Parser     format.Format
	UnixSocket *string
	TCPService *string
	UDPService *string
}

func (s *ServerOptions) String() string {
	return "{\n" +
		fmt.Sprintf("    TCPService: %s\n", *s.TCPService) +
		fmt.Sprintf("    UDPService: %s\n", *s.UDPService) +
		fmt.Sprintf("    UnixSocket: %s\n", *s.UnixSocket) +
		"  }"
}

// IsPassThrough returns true when skipping parsing
func (s *ServerOptions) IsPassThrough() bool {
	return s.Parser == domain.RAW
}

type consoleOpts struct {
	Trace  *bool
	Silent *bool
}

func (c *consoleOpts) Bind(opts *Options) {
	*c.Trace = *opts.Trace
	*c.Silent = *opts.Silent
}

func (c *consoleOpts) Invalid() bool {
	return *c.Trace && *c.Silent
}

// flag.Bool("cat", false, "Forward the input data as-is; do not try to parse as syslog-standard formats")
