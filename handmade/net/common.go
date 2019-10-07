package net

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

type (
	Command string

	Id   string
	Time int64

	EventCallback func(event Event)

	WebsocketClienter interface {
		DialAndServe(ctx context.Context, urlStr string) error

		// SendMessage делает копию data, т.е. data можно модифицировать сразу после вызова
		SendMessage(event Event, cb EventCallback) error

		IsConnected() bool

		SetEventsReceiver(EventCallback)
	}
)

var (
	ErrWSIsNotConnected = fmt.Errorf(`WS is not connected`)

	UnknownPlayer = Id(``)
)

var (
	// базовые
	CommandEmpty = Command(``)
	CommandPing  = Command(`ping`)
	CommandPong  = Command(`pong`)
	CommandError = Command(`error`)

	// взаимодействие сторон
	CommandJoin        = Command(`join`)
	CommandRaceInfo    = Command(`race_info`)
	CommandRaceState   = Command(`race_state`)
	CommandRaceNew     = Command(`race_new`)
	CommandIAmReady    = Command(`i_am_ready`)
	CommandIAmFinished = Command(`i_am_finished`)
	CommandRaceStart   = Command(`race_start`)
	CommandPlayerStep  = Command(`player_step`)
)

var (
	CommandErrorUnknownCommand = `unknown_command`
)

func NewCommandError(data string) Event {
	return NewEvent(CommandError, data)
}

func NewCommandEmpty() Event {
	return NewEvent(CommandEmpty, ``)
}

func GenId() Id {
	var randBuf [16]byte // 128бит хватит всем :)
	rand.Read(randBuf[:])
	return Id(hex.EncodeToString(randBuf[:]))
}

func GetCurrentTime() Time {
	return Time(time.Now().UnixNano() / int64(time.Millisecond))
}
