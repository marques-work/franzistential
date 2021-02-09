package conf

import (
	"fmt"
	"strings"

	"github.com/marques-work/franzistential/domain"
	"gopkg.in/mcuadros/go-syslog.v2/format"
)

var Format format.Format

var (
	RAW     format.Format = &domain.RawFormat{}
	RFC3164 format.Format = &format.RFC3164{}
	RFC5424 format.Format = &format.RFC5424{}
	RFC6587 format.Format = &format.RFC6587{}
	DETECT  format.Format = &format.Automatic{}
)

type FormatFlag struct{}

func (f *FormatFlag) Set(mode string) error {
	switch strings.ToLower(mode) {
	case "raw":
		Format = RAW
	case "guess-rfc":
		Format = DETECT
	case "rfc3164":
		Format = RFC3164
	case "rfc5424":
		Format = RFC5424
	case "rfc6587":
		Format = RFC6587
	default:
		return fmt.Errorf("Unknown format `%s`; valid formats are [ raw | guess-rfc | rfc3164 | rfc5424 | rfc6587 ]", mode)
	}

	return nil
}
