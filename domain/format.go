package domain

import (
	"bufio"
	"time"

	"gopkg.in/mcuadros/go-syslog.v2/format"
)

var (
	// RAW represents a no-op/pass-through syslog parser
	RAW format.Format = &rawFormat{}

	// RFC3164 syslog parser
	RFC3164 format.Format = &format.RFC3164{}

	// RFC5424 syslog parser
	RFC5424 format.Format = &format.RFC5424{}

	// RFC6587 syslog parser
	RFC6587 format.Format = &format.RFC6587{}

	// DETECT represents a syslog parser that does a best-guess in parsing RFCXXXX syslog formats
	DETECT format.Format = &format.Automatic{}
)

type rawFormat struct {
}

func (f *rawFormat) GetParser(payload []byte) format.LogParser {
	return newPassThruParser(payload)
}

func (f *rawFormat) GetSplitFunc() bufio.SplitFunc {
	return bufio.ScanLines
}

func newPassThruParser(payload []byte) *passThruParser {
	p := make(format.LogParts)

	message := string(payload)
	p["rawmsg"] = message
	p["rawmsg-after-pri"] = message
	p["msg"] = message
	return &passThruParser{Parts: p}
}

type passThruParser struct {
	Parts format.LogParts
}

func (p *passThruParser) Dump() format.LogParts {
	return p.Parts
}

func (p *passThruParser) Parse() error {
	return nil
}

func (p *passThruParser) Location(t *time.Location) {
}
