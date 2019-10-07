package events

import "errors"

const (
	floatFixedMult        = 1000
	formatSeparator       = ':'
	raceStateFieldsCount  = 4 // количество сериализуемых полей в PlayerCompactState
	playerStepFieldsCount = 5 // количество сериализуемых полей в EventPlayerStep
)

var (
	ErrWrongFormat = errors.New(`wrong format`)
)
