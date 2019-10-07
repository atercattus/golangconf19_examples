package events

import (
	"encoding/json"
	"github.com/atercattus/golangconf19_examples/engo/net"
)

type (
	EventUsersCount struct {
		Count int32
	}
)

var (
	_ NetEventer = &EventUsersCount{}
)

func (event *EventUsersCount) ToEvent() net.Event {
	return net.NewEvent(
		net.CommandUsersCount,
		string(event.Marshal()),
	)
}

func (event *EventUsersCount) Marshal() []byte {
	data, _ := json.Marshal(event)
	return data
}

func (event *EventUsersCount) Unmarshal(data []byte) error {
	return json.Unmarshal(data, event)
}
