package events

import (
	"encoding/json"
	"github.com/atercattus/golangconf19_examples/handmade/net"
)

type (
	PlayerInfo struct {
		Name       string
		Id         net.Id
		Distance   float32
		Speed      float32
		Ready      bool
		FinishedAt net.Time
	}

	PlayerInfoWithTimes struct {
		PlayerInfo

		LastStepAtServer net.Time
		LastStepAtPlayer net.Time
	}

	EventRaceInfo struct {
		EventBase

		Players []PlayerInfo
	}
)

func (event *EventRaceInfo) getCommand() net.Command {
	return net.CommandRaceInfo
}

func (event *EventRaceInfo) Marshal() []byte {
	data, _ := json.Marshal(event)
	return data
}

func (event *EventRaceInfo) Unmarshal(data []byte) error {
	return json.Unmarshal(data, event)
}
