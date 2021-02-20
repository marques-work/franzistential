package conf

import (
	"errors"
	"fmt"
	"os"

	"github.com/marques-work/franzistential/domain"
	"github.com/marques-work/franzistential/logging"
	"gopkg.in/mcuadros/go-syslog.v2/format"
)

// Instance is the default Options value
var Instance = &Options{}

// Configure initializes and validates the default Options
func Configure() {
	if err := Instance.Validate(); err != nil {
		logging.Die("INVALID CONFIGURATION: %s", err.Error())
	}

	Instance.Bind(&consoleBinding{
		Trace:  &logging.Trace,
		Silent: &logging.Silent,
	})
}

// Options represents configuration
type Options struct {
	Out          *bool // keeping STDOUT as a boolean ensures it can only be specified once
	Destinations []domain.Destination
	SendTimeout  *uint64

	// logging
	Silent *bool
	Trace  *bool

	Server *ServerOptions

	_parsed bool
}

// ServerMode returns whether or not to start a daemon
func (o *Options) ServerMode() bool {
	return nil != o.Server
}

// HasDestination returns whether or not there is at least one output destination configured
func (o *Options) HasDestination() bool {
	return *o.Out || len(o.Destinations) > 0
}

// Bind parses and applies options to the target
func (o *Options) Bind(targetConsole *consoleBinding) {
	targetConsole.Bind(o)

	if !o._parsed /* this should ever only happen once */ && *o.Out {
		d, _ := domain.NewIO(os.Stdout)
		o.Destinations = append(o.Destinations, d)
	}

	o._parsed = true
}

func (o *Options) Validate() error {
	if *o.Trace && *o.Silent {
		return errors.New("Both `--verbose` and `--quiet` cannot be simultaneously set; pick one")
	}

	if !o.HasDestination() {
		logging.Warn("You have not configured any outputs; this is effectively a black hole, and might be a misconfiguration.")
	}

	if o.ServerMode() {
		tcp := *o.Server.TCPService

		if "" != tcp {
			if err := validateNetworkListenerConfig("tcp", tcp); err != nil {
				return err
			}
		}

		udp := *o.Server.UDPService

		if "" != udp {
			if err := validateNetworkListenerConfig("udp", udp); err != nil {
				return err
			}
		}
	}

	return nil
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
		fmt.Sprintf("    Parser: %T\n", s.Parser) +
		fmt.Sprintf("    TCPService: %s\n", *s.TCPService) +
		fmt.Sprintf("    UDPService: %s\n", *s.UDPService) +
		fmt.Sprintf("    UnixSocket: %s\n", *s.UnixSocket) +
		"  }"
}

// IsPassThrough returns true when skipping parsing
func (s *ServerOptions) IsPassThrough() bool {
	return s.Parser == RAW
}

type consoleBinding struct {
	Trace  *bool
	Silent *bool
}

func (c *consoleBinding) Bind(opts *Options) {
	*c.Trace = *opts.Trace
	*c.Silent = *opts.Silent
}
