package events

import (
	"bytes"
	"github.com/atercattus/golangconf19_examples/handmade/net"
	"strconv"
)

type (
	// Это очень частое событие. Чтобы немного сэкономить на трафике, делаю свой упоротый мини формат
	EventPlayerStep struct {
		EventBase

		PlayerId net.Id
		RaceId   net.Id
		Distance float32
		Speed    float32
		Now      net.Time
	}
)

func (event *EventPlayerStep) getCommand() net.Command {
	return net.CommandPlayerStep
}

func (event *EventPlayerStep) Marshal() []byte {
	var data []byte

	data = strconv.AppendInt(data, int64(event.Distance*floatFixedMult), 10)
	data = append(data, formatSeparator)
	data = strconv.AppendInt(data, int64(event.Speed*floatFixedMult), 10)
	data = append(data, formatSeparator)
	data = strconv.AppendInt(data, int64(event.Now), 10)
	data = append(data, formatSeparator)
	data = append(data, event.PlayerId...)
	data = append(data, formatSeparator) // hex не будет конфликтовать с разделителем
	data = append(data, event.RaceId...)

	return data
}

func (event *EventPlayerStep) Unmarshal(data []byte) error {
	chunks := bytes.Split(data, []byte{formatSeparator})
	if len(chunks) != playerStepFieldsCount {
		return ErrWrongFormat
	}

	event.PlayerId = net.Id(chunks[3])
	event.RaceId = net.Id(chunks[4])

	chunks = chunks[:3]
	ints := make([]int64, len(chunks))
	for i, chunk := range chunks {
		if val, err := strconv.ParseInt(string(chunk), 10, 64); err != nil {
			return ErrWrongFormat
		} else {
			ints[i] = val
		}
	}

	event.Distance = float32(ints[0]) / floatFixedMult
	event.Speed = float32(ints[1]) / floatFixedMult
	event.Now = net.Time(ints[2])

	return nil
}
