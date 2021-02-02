package domain

import (
	"bufio"
	"time"

	"gopkg.in/mcuadros/go-syslog.v2/format"
)

type RawFormat struct {
}

func (f *RawFormat) GetParser(payload []byte) format.LogParser {
	return NewNoopParser(payload)
}

func (f *RawFormat) GetSplitFunc() bufio.SplitFunc {
	return bufio.ScanLines
}

func NewNoopParser(payload []byte) *NoopParser {
	var p format.LogParts
	message := string(payload)
	p["rawmsg"] = message
	p["rawmsg-after-pri"] = message
	p["msg"] = message
	return &NoopParser{Parts: p}
}

type NoopParser struct {
	Parts format.LogParts
}

func (p *NoopParser) Dump() format.LogParts {
	return p.Parts
}

func (p *NoopParser) Parse() error {
	return nil
}

func (p *NoopParser) Location(t *time.Location) {
}
