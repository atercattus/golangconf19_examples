package net

import (
	"encoding/json"
	"log"
	"sync/atomic"
)

type (
	EventId int32

	Event struct {
		Idx     EventId
		Command Command
		Data    string
	}
)

var (
	lastEventIdx int32
)

func NewEvent(command Command, data string) Event {
	return Event{
		Idx:     EventId(atomic.AddInt32(&lastEventIdx, 1)),
		Command: command,
		Data:    data,
	}
}

func (ev *Event) Marshal() []byte {
	data, err := json.Marshal(ev)
	if err != nil {
		log.Println("JSON marshall error:", err)
	}
	return data
}

func (ev *Event) Unmarshal(data []byte) error {
	return json.Unmarshal(data, ev)
}
