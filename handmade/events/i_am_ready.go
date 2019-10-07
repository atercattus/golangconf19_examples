package events

import (
	"encoding/json"
	"github.com/atercattus/golangconf19_examples/handmade/net"
)

type (
	EventIAmReady struct {
		EventBase

		PlayerId net.Id
		RaceId   net.Id
	}
)

func (event *EventIAmReady) getCommand() net.Command {
	return net.CommandIAmReady
}

func (event *EventIAmReady) Marshal() []byte {
	data, _ := json.Marshal(event)
	return data
}

func (event *EventIAmReady) Unmarshal(data []byte) error {
	return json.Unmarshal(data, event)
}
