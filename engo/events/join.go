package events

import (
	"encoding/json"
	"github.com/atercattus/golangconf19_examples/engo/net"
)

type (
	EventJoin struct {
		UserName string
	}
)

var (
	_ NetEventer = &EventJoin{}
)

func (event *EventJoin) ToEvent() net.Event {
	return net.NewEvent(
		net.CommandJoin,
		string(event.Marshal()),
	)
}

func (event *EventJoin) Marshal() []byte {
	data, _ := json.Marshal(event)
	return data
}

func (event *EventJoin) Unmarshal(data []byte) error {
	return json.Unmarshal(data, event)
}
