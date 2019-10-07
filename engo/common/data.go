package common

import (
	"github.com/atercattus/golangconf19_examples/engo/events"
	"github.com/atercattus/golangconf19_examples/engo/net"
)

type (
	GlobalDataState struct {
		UserName   string
		UsersCount int32
		WS         net.WebsocketClienter
	}
)

var (
	GlobalData GlobalDataState
)

func (gd *GlobalDataState) Join(username string, cb net.EventCallback) error {
	gd.UserName = username

	var ev events.EventJoin
	ev.UserName = username

	return gd.WS.SendMessage(ev.ToEvent(), cb)
}
