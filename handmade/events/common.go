package events

import (
	"github.com/atercattus/golangconf19_examples/handmade/net"
)

type (
	Eventer interface {
		getCommand() net.Command
		Marshal() []byte
		Unmarshal(data []byte) error
	}

	EventBase struct {
	}
)

var (
	_ Eventer = &EventBase{}
)

func EventToNetEvent(event Eventer) net.Event {
	return net.NewEvent(
		event.getCommand(),
		string(event.Marshal()),
	)
}

func (event *EventBase) getCommand() net.Command {
	return net.CommandEmpty
}

func (event *EventBase) Marshal() []byte {
	return nil
}

func (event *EventBase) Unmarshal(data []byte) error {
	return nil
}
