package events

import "github.com/atercattus/golangconf19_examples/engo/net"

type (
	NetEventer interface {
		ToEvent() net.Event
	}
)
