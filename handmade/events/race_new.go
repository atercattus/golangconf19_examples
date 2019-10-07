package events

import (
	"encoding/json"
	"github.com/atercattus/golangconf19_examples/handmade/net"
)

type (
	EventRaceNew struct {
		EventBase

		RaceUuid   net.Id
		MaxPlayers int // для UI
		Distance   float32
	}
)

func (event *EventRaceNew) getCommand() net.Command {
	return net.CommandRaceNew
}

func (event *EventRaceNew) Marshal() []byte {
	data, _ := json.Marshal(event)
	return data
}

func (event *EventRaceNew) Unmarshal(data []byte) error {
	return json.Unmarshal(data, event)
}
