package net

import (
	"context"
	"fmt"
)

type (
	Command string

	EventCallback func(event Event)

	WebsocketClienter interface {
		DialAndServe(ctx context.Context, urlStr string) error

		// SendMessage делает копию data, т.е. data можно модифицировать сразу после вызова
		SendMessage(event Event, cb EventCallback) error

		IsConnected() bool
	}
)

var (
	ErrWSIsNotConnected = fmt.Errorf(`WS is not connected`)
)

var (
	CommandEmpty      = Command(``)
	CommandPing       = Command(`ping`)
	CommandPong       = Command(`pong`)
	CommandError      = Command(`error`)
	CommandJoin       = Command(`join`)
	CommandUsersCount = Command(`users_count`)
)

var (
	CommandErrorUnknownCommand = `unknown_command`
)

func NewCommandError(data string) Event {
	return NewEvent(CommandError, data)
}
