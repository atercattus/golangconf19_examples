package events

import (
	"encoding/json"
	"github.com/atercattus/golangconf19_examples/handmade/net"
)

type (
	EventRaceStart struct {
		EventBase

		RaceId net.Id
	}
)

func (event *EventRaceStart) getCommand() net.Command {
	return net.CommandRaceStart
}

func (event *EventRaceStart) Marshal() []byte {
	data, _ := json.Marshal(event)
	return data
}

func (event *EventRaceStart) Unmarshal(data []byte) error {
	return json.Unmarshal(data, event)
}
