package cli

import (
	"fmt"
	"strings"

	"github.com/marques-work/franzistential/conf"
	"github.com/marques-work/franzistential/domain"
	"github.com/spf13/pflag"
	"gopkg.in/mcuadros/go-syslog.v2/format"
)

var (
	knownParsers = map[string]format.Format{
		"raw":       conf.RAW,
		"guess-rfc": conf.DETECT,
		"rfc3164":   conf.RFC3164,
		"rfc5424":   conf.RFC5424,
		"rfc6587":   conf.RFC6587,
	}
)

func newEventHubValue(opts *conf.Options) pflag.Value {
	return &eventHubValue{opts: opts}
}

type eventHubValue struct {
	raw  string
	opts *conf.Options
}

func (eh *eventHubValue) Type() string {
	return "string"
}

func (eh *eventHubValue) Set(url string) error {
	if eh.raw == url {
		return nil
	}

	eh.raw = url
	dest, err := domain.NewEventHub(url)

	if err == nil {
		eh.opts.Destinations = append(eh.opts.Destinations, dest)
	}

	return err
}

func (eh *eventHubValue) String() string { return eh.raw }

func newParserValue(opts *conf.ServerOptions) pflag.Value {
	return &parserValue{opts: opts}
}

type parserValue struct {
	raw  string
	opts *conf.ServerOptions
}

func (f *parserValue) Type() string {
	return "string"
}

func (f *parserValue) Set(mode string) error {
	f.raw = strings.ToLower(mode)

	if p, ok := knownParsers[f.raw]; ok {
		f.opts.Parser = p
	} else {
		return fmt.Errorf("Unknown parser `%s`; valid parsers are [ %s ]", mode, parserList())
	}

	return nil
}

func (f *parserValue) String() string {
	return f.raw
}

func parserList() string {
	keys := make([]string, 0, len(knownParsers))
	for k := range knownParsers {
		keys = append(keys, k)
	}

	return strings.Join(keys, " | ")
}
