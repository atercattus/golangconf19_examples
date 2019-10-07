package events

import (
	"bytes"
	"github.com/atercattus/golangconf19_examples/handmade/net"
	"strconv"
)

type (
	PlayerCompactState struct {
		Distance         float32
		Speed            float32
		LastStepAtServer net.Time
		FinishedAt       net.Time
	}

	// Это очень частое событие. Чтобы немного сэкономить на трафике, делаю свой упоротый мини формат
	EventRaceState struct {
		EventBase

		Players []PlayerCompactState
	}
)

func (event *EventRaceState) getCommand() net.Command {
	return net.CommandRaceState
}

func (event *EventRaceState) Marshal() []byte {
	var data []byte

	data = strconv.AppendInt(data, int64(len(event.Players)), 10)
	data = append(data, formatSeparator)
	for _, player := range event.Players {
		data = strconv.AppendInt(data, int64(player.Distance*floatFixedMult), 10)
		data = append(data, formatSeparator)
		data = strconv.AppendInt(data, int64(player.Speed*floatFixedMult), 10)
		data = append(data, formatSeparator)
		data = strconv.AppendInt(data, int64(player.LastStepAtServer), 10)
		data = append(data, formatSeparator)
		data = strconv.AppendInt(data, int64(player.FinishedAt), 10)
		data = append(data, formatSeparator)
	}
	data = bytes.TrimSuffix(data, []byte{formatSeparator})

	return data
}

func (event *EventRaceState) Unmarshal(data []byte) error {
	chunks := bytes.Split(data, []byte{formatSeparator})
	ints := make([]int64, len(chunks))
	for i, chunk := range chunks {
		if val, err := strconv.ParseInt(string(chunk), 10, 64); err != nil {
			return ErrWrongFormat
		} else {
			ints[i] = val
		}
	}

	cnt := int(ints[0])
	ints = ints[1:]
	if cnt*raceStateFieldsCount != len(ints) {
		return ErrWrongFormat
	}

	event.Players = make([]PlayerCompactState, cnt)
	for i := 0; i < cnt; i++ {
		player := &event.Players[i]
		player.Distance = float32(ints[raceStateFieldsCount*i]) / floatFixedMult
		player.Speed = float32(ints[raceStateFieldsCount*i+1]) / floatFixedMult
		player.LastStepAtServer = net.Time(ints[raceStateFieldsCount*i+2])
		player.FinishedAt = net.Time(ints[raceStateFieldsCount*i+3])
	}
	return nil
}
