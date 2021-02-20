package conf_test

import (
	"testing"

	"github.com/marques-work/franzistential/conf"
)

func TestPassThruParser(t *testing.T) {
	as := asserts(t)

	parser := conf.RAW.GetParser([]byte("hello"))

	as.ok(parser.Parse())
	as.eq("hello", parser.Dump()["msg"])
}
