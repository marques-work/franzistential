package cli

import (
	"fmt"
	"strings"

	"github.com/marques-work/franzistential/conf"
	"github.com/marques-work/franzistential/domain"
	"gopkg.in/mcuadros/go-syslog.v2/format"

	"github.com/spf13/pflag"
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
	if dest, err := domain.NewEventHub(url); err == nil {
		eh.opts.Destinations = append(eh.opts.Destinations, dest)
		return nil
	} else {
		return err
	}
}

func (eh *eventHubValue) String() string { return eh.raw }

type FormatFlag struct {
	Format format.Format
}

func (f *FormatFlag) Set(mode string) error {
	switch strings.ToLower(mode) {
	case "raw":
		f.Format = RAW
	case "guess-rfc":
		f.Format = DETECT
	case "rfc3164":
		f.Format = RFC3164
	case "rfc5424":
		f.Format = RFC5424
	case "rfc6587":
		f.Format = RFC6587
	default:
		return fmt.Errorf("Unknown format `%s`; valid formats are [ raw | guess-rfc | rfc3164 | rfc5424 | rfc6587 ]", mode)
	}

	return nil
}

var (
	RAW     format.Format = &domain.RawFormat{}
	RFC3164 format.Format = &format.RFC3164{}
	RFC5424 format.Format = &format.RFC5424{}
	RFC6587 format.Format = &format.RFC6587{}
	DETECT  format.Format = &format.Automatic{}
)
