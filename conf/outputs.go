package conf

import "github.com/marques-work/franzistential/domain"

type EventHubFlag struct{}

func (eh *EventHubFlag) Set(url string) error {
	if dest, err := domain.NewEventHub(url); err == nil {
		Destinations = append(Destinations, dest)
		return nil
	} else {
		return err
	}
}

// I don't think we care about the string rep since it's just a shim
func (eh *EventHubFlag) String() string { return "-eventHub" }

type UnixgramOutFlag struct{}
