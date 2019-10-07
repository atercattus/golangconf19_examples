package main

import (
	"fmt"
	"github.com/atercattus/golangconf19_examples/handmade/events"
	"github.com/atercattus/golangconf19_examples/handmade/net"
	"log"
	"sync"
	"time"
)

type (
	GameState struct {
		lock       sync.Mutex
		futureRace *Race
		races      map[net.Id]*Race
		wsServer   *WSServer
	}
)

const (
	// Длина трассы (метры к примеру). ToDo: неплохо бы сделать вариативность (короткие, средние, марафонские)
	raceDistance           = 2000 // ToDo: поменять на 3000-5000 ?
	raceBeforeTickInterval = 1 * time.Second
	raceInGameTickInterval = 500 * time.Millisecond
)

var (
	ErrUnknownRaceId   = fmt.Errorf(`unknown raceId`)
	ErrUnknownPlayerId = fmt.Errorf(`unknown playerId`)
)

func NewGameState(server *WSServer) *GameState {
	return &GameState{
		races:    make(map[net.Id]*Race),
		wsServer: server,
	}
}

func (gs *GameState) JoinPlayer(ev events.EventJoin) (*Race, error) {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	race := gs.getOrCreateFutureRace()
	if race.Join(ev) {
		gs.futureRace = nil
	}

	return race, nil
}

func (gs *GameState) IAmReady(ev events.EventIAmReady) error {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	if race, ok := gs.races[ev.RaceId]; !ok {
		return ErrUnknownRaceId
	} else if !race.IAmReady(ev.PlayerId) {
		return ErrUnknownPlayerId
	}
	return nil
}

func (gs *GameState) IAmFinished(ev events.EventIAmFinished) error {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	if race, ok := gs.races[ev.RaceId]; !ok {
		return ErrUnknownRaceId
	} else if !race.IAmFinished(ev.PlayerId) {
		return ErrUnknownPlayerId
	}
	return nil
}

func (gs *GameState) PlayerStep(ev events.EventPlayerStep) error {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	if race, ok := gs.races[ev.RaceId]; !ok {
		return ErrUnknownRaceId
	} else if !race.UpdatePlayerPos(ev) {
		return ErrUnknownPlayerId
	}
	return nil
}

func (gs *GameState) getFutureRace() *Race {
	if (gs.futureRace != nil) && !gs.futureRace.StartTime.IsZero() {
		// гонка уже началась
		gs.futureRace = nil
	}
	return gs.futureRace
}

func (gs *GameState) getOrCreateFutureRace() *Race {
	if race := gs.getFutureRace(); race != nil {
		return race
	}

	gs.futureRace = NewRace(raceDistance)
	gs.races[gs.futureRace.Id] = gs.futureRace
	go gs.raceProcessor(gs.futureRace)

	return gs.futureRace
}

func (gs *GameState) raceProcessorBeforeRaceStart(race *Race) bool {
	ticker := time.NewTicker(raceBeforeTickInterval)
	defer ticker.Stop()

	var evRaceInfo events.EventRaceInfo

	for range ticker.C {
		// уг. в идеале под локом целиком сериализовать
		evRaceInfo.Players = race.GetPlayers()
		if !gs.wsServer.BroadcastEventToRace(race.Id, events.EventToNetEvent(&evRaceInfo)) {
			// все отключились
			println(`Race`, race.Id, `send EventRaceInfo failed. All are disconnected`)
			return false
		}

		if readyCnt := race.GetReadyCount(); readyCnt < 2 {
			// ToDo: если есть еще участники, то это странно (двое одобрили - стартуют все)
			continue
		}

		// пора создавать гонку
		race.Start()
		var evRaceStart events.EventRaceStart
		evRaceStart.RaceId = race.Id
		if !gs.wsServer.BroadcastEventToRace(race.Id, events.EventToNetEvent(&evRaceStart)) {
			// все отключились
			println(`Race`, race.Id, `send EventRaceStart failed. All are disconnected`)
			return false
		}

		break
	}

	return true
}

func (gs *GameState) raceProcessorAfterRaceStart(race *Race) bool {
	ticker := time.NewTicker(raceInGameTickInterval)
	defer ticker.Stop()

	var evRaceState events.EventRaceState

	for range ticker.C {
		evRaceState.Players = race.GetPlayersState()
		if !gs.wsServer.BroadcastEventToRace(race.Id, events.EventToNetEvent(&evRaceState)) {
			// все отключились
			println(`Race`, race.Id, `send EventRaceState failed. All are disconnected`)
			return false
		}

		// ToDo: Крутим лупы, пока есть коннекты. Нужно условие завершения
	}

	return true
}

func (gs *GameState) raceProcessor(race *Race) {
	defer log.Println(`GS race done`, race.Id)

	if !gs.raceProcessorBeforeRaceStart(race) {
		return
	}
	gs.raceProcessorAfterRaceStart(race)
}

func (gs *GameState) DeletePlayerFromFutureRace(playerId net.Id) {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	if race := gs.getFutureRace(); race != nil {
		if isEmpty := race.Leave(playerId); isEmpty {
			// в гонке никого не осталось, удаляю ее
			delete(gs.races, race.Id)
			gs.futureRace = nil
		}
	}
}
