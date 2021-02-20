package domain_test

import (
	"os"
	"testing"

	"github.com/marques-work/franzistential/domain"
)

func TestEventHubRejectsBadConnectionString(t *testing.T) {
	as := asserts(t)
	h, err := domain.NewEventHub("definitely not a url")

	as.err("Malformed Event Hub connection string: failed parsing connection string due to unmatched key value separated by '='", err)
	as.isNil(h)
}

func TestEventHubParses(t *testing.T) {
	as := asserts(t)
	h, err := domain.NewEventHub("Endpoint=sb://test.a.b.c/;SharedAccessKeyName=Foo;SharedAccessKey=secret;EntityPath=hubname")

	as.ok(err)
	as.isNotNil(h)

	as.eq("Endpoint=sb://test.a.b.c/;SharedAccessKeyName=Foo;SharedAccessKey=secret;EntityPath=hubname", h.ConnectString())
	as.eq("Endpoint=sb://test.a.b.c/;SharedAccessKeyName=Foo;SharedAccessKey=********;EntityPath=hubname", h.Redacted())
}

func TestIODest(t *testing.T) {
	as := asserts(t)
	i, err := domain.NewIO(os.Stderr)

	as.ok(err)
	as.isNotNil(i)

	as.eq(os.Stderr.Name(), i.ConnectString())
	as.eq(i.ConnectString(), i.Redacted())
}
