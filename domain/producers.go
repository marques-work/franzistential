package domain

import (
	"fmt"
	"io"
	"time"

	"github.com/marques-work/franzistential/util"
)

type Producer interface {
	StartForwarding() error
	Accept(line string) error
	Teardown() error
}

type eventHubProducer struct {
	hub   *EventHub
	input chan string   // data in
	queue chan []string // data out

	_running bool
	_closed  bool
}

func (e *eventHubProducer) StartForwarding() error {
	return nil
}

func (e *eventHubProducer) Accept(line string) error {
	return nil
}

func (e *eventHubProducer) Teardown() error {
	if !e._running {
		return fmt.Errorf("Cannot teardown an EventHubProducer that is not running")
	}

	if e._closed {
		return fmt.Errorf("Event Hub Producer cannot be torn down twice")
	}

	e._closed = true
	defer close(e.queue)
	return nil
}

func NewEventHubProducer(connect string, data chan string) (Producer, error) {
	hub, err := NewEventHub(connect)

	if err != nil {
		return nil, err
	}

	return &eventHubProducer{
		hub:   hub,
		input: data,
		queue: util.SpaceTimeSlicer(data, 1024, time.Duration(500)*time.Millisecond),
	}, nil
}

// simple, as opposed to batched and/or buffered; perfect for non-disk streams, like STDOUT
// which can handle immediate writes
type simpleIOProducer struct {
	io.Writer
}
