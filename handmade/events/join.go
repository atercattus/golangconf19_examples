package events

import (
	"encoding/json"
	"github.com/atercattus/golangconf19_examples/handmade/net"
)

type (
	EventJoin struct {
		EventBase

		PlayerName string
		PlayerId   net.Id
	}
)

func (event *EventJoin) getCommand() net.Command {
	return net.CommandJoin
}

func (event *EventJoin) Marshal() []byte {
	data, _ := json.Marshal(event)
	return data
}

func (event *EventJoin) Unmarshal(data []byte) error {
	return json.Unmarshal(data, event)
}
