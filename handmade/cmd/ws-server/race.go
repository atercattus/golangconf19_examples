package main

import (
	"github.com/atercattus/golangconf19_examples/handmade/events"
	"github.com/atercattus/golangconf19_examples/handmade/net"
	"log"
	"sort"
	"strings"
	"sync"
	"time"
)

type (
	Race struct {
		lock sync.Mutex

		Id          net.Id
		players     []events.PlayerInfoWithTimes
		Distance    float32
		CreatedTime time.Time
		StartTime   time.Time
	}
)

const (
	maxPlayersInRace = 4
	//maxTimeForWaitingGroup = 20 * time.Second
)

func NewRace(distance float32) *Race {
	race := &Race{
		Id:          net.GenId(),
		Distance:    distance,
		CreatedTime: time.Now(),
	}
	return race
}

func (race *Race) Join(ev events.EventJoin) (full bool) {
	race.lock.Lock()
	defer race.lock.Unlock()

	if !race.StartTime.IsZero() {
		return false
	}

	var playerInfo events.PlayerInfoWithTimes
	playerInfo.Name = ev.PlayerName
	playerInfo.Id = ev.PlayerId
	race.players = append(race.players, playerInfo)
	full = len(race.players) == maxPlayersInRace
	race.sortPlayers()

	return full
}

func (race *Race) Leave(playerId net.Id) (empty bool) {
	race.lock.Lock()
	defer race.lock.Unlock()

	if !race.StartTime.IsZero() {
		return false
	}

	pos := -1
	for p, player := range race.players {
		if player.Id == playerId {
			pos = p
			break
		}
	}

	if pos > -1 {
		last := len(race.players) - 1
		if pos < last {
			race.players[pos] = race.players[last]
		}
		race.players = race.players[:last]
		race.sortPlayers()
	}

	return len(race.players) == 0
}

func (race *Race) IAmReady(playerId net.Id) (succ bool) {
	race.lock.Lock()
	defer race.lock.Unlock()

	for idx := range race.players {
		player := &race.players[idx]
		if player.Id == playerId {
			player.Ready = true
			return true
		}
	}
	return false
}

func (race *Race) IAmFinished(playerId net.Id) (succ bool) {
	race.lock.Lock()
	defer race.lock.Unlock()

	for idx := range race.players {
		player := &race.players[idx]
		if player.Id == playerId && player.FinishedAt == 0 {
			player.FinishedAt = net.GetCurrentTime()
			log.Printf(`Player %s (%s) finished in race %s`, player.Name, player.Id, race.Id)
			return true
		}
	}
	return false
}

func (race *Race) UpdatePlayerPos(ev events.EventPlayerStep) (succ bool) {
	race.lock.Lock()
	defer race.lock.Unlock()

	for idx := range race.players {
		player := &race.players[idx]
		if player.Id == ev.PlayerId {
			player.Distance = ev.Distance
			player.Speed = ev.Speed
			player.LastStepAtPlayer = ev.Now
			player.LastStepAtServer = net.GetCurrentTime()
			return true
		}
	}
	return false
}

func (race *Race) GetPlayers() []events.PlayerInfo {
	race.lock.Lock()
	defer race.lock.Unlock()

	lst := make([]events.PlayerInfo, len(race.players))
	for i := range race.players {
		lst[i] = race.players[i].PlayerInfo
	}
	return lst
}

func (race *Race) GetPlayersState() []events.PlayerCompactState {
	race.lock.Lock()
	defer race.lock.Unlock()

	lst := make([]events.PlayerCompactState, len(race.players))
	for i := range race.players {
		player := &race.players[i]
		lst[i].Distance = player.Distance
		lst[i].Speed = player.Speed
		lst[i].LastStepAtServer = player.LastStepAtServer
		lst[i].FinishedAt = player.FinishedAt
	}
	return lst
}

func (race *Race) Start() bool {
	race.lock.Lock()
	defer race.lock.Unlock()

	succ := race.StartTime.IsZero()
	if succ {
		race.StartTime = time.Now()
	}
	return succ
}

func (race *Race) GetReadyCount() int {
	race.lock.Lock()
	defer race.lock.Unlock()

	cnt := 0
	for i := range race.players {
		if race.players[i].Ready {
			cnt++
		}
	}
	return cnt
}

func (race *Race) sortPlayers() {
	sort.Slice(race.players, func(i, j int) bool {
		return strings.Compare(race.players[i].Name, race.players[j].Name) < 0
	})
}
