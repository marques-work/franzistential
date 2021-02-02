package conf

import "github.com/marques-work/franzistential/domain"

type EventHubFlag struct{}

func (eh *EventHubFlag) Set(url string) error {
	Destinations = append(Destinations, domain.NewEventHub(url))
	return nil
}

// I don't think we care about the string rep since it's just a shim
func (eh *EventHubFlag) String() string { return "[event-hub-uri]" }
