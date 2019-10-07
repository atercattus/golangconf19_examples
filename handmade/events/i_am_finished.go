package events

import (
	"encoding/json"
	"github.com/atercattus/golangconf19_examples/handmade/net"
)

type (
	EventIAmFinished struct {
		EventBase

		PlayerId net.Id
		RaceId   net.Id
	}
)

func (event *EventIAmFinished) getCommand() net.Command {
	return net.CommandIAmFinished
}

func (event *EventIAmFinished) Marshal() []byte {
	data, _ := json.Marshal(event)
	return data
}

func (event *EventIAmFinished) Unmarshal(data []byte) error {
	return json.Unmarshal(data, event)
}
