package main

import (
	"github.com/atercattus/golangconf19_examples/handmade/net"
	"sync"
)

type (
	EventChan chan *net.Event

	ChanToPlayerInfo struct {
		PlayerId net.Id
		RaceId   net.Id
	}

	RaceIdToPlayerConnections struct {
		sync.RWMutex
		Mapping     map[net.Id][]EventChan
		BackMapping map[EventChan]ChanToPlayerInfo
	}
)

func NewRaceIdToPlayerConnections() *RaceIdToPlayerConnections {
	m := &RaceIdToPlayerConnections{
		Mapping:     make(map[net.Id][]EventChan),
		BackMapping: make(map[EventChan]ChanToPlayerInfo),
	}
	return m
}

func (m *RaceIdToPlayerConnections) AddConnection(raceId net.Id, playerId net.Id, sendChan EventChan) {
	m.Lock()
	lst := m.Mapping[raceId]
	lst = append(lst, sendChan)
	m.Mapping[raceId] = lst
	m.BackMapping[sendChan] = ChanToPlayerInfo{PlayerId: playerId, RaceId: raceId}
	m.Unlock()
}

func (m *RaceIdToPlayerConnections) GetChans(raceId net.Id) []EventChan {
	m.RLock()
	lst := m.Mapping[raceId]
	m.RUnlock()
	return lst
}

func (m *RaceIdToPlayerConnections) deleteChanInList(sendChan EventChan, lst []EventChan) []EventChan {
	pos := -1
	for p, ch := range lst {
		if ch == sendChan {
			pos = p
			break
		}
	}

	if pos > -1 {
		last := len(lst) - 1
		if pos < last {
			lst[pos] = lst[last]
			lst[last] = nil
		}
		lst = lst[:last]
	}

	return lst
}

func (m *RaceIdToPlayerConnections) DeleteConnectionBySendChan(sendChan EventChan) (playerId net.Id) {
	m.Lock()
	defer m.Unlock()

	if info, ok := m.BackMapping[sendChan]; !ok {
		playerId = net.UnknownPlayer
	} else {
		playerId = info.PlayerId

		delete(m.BackMapping, sendChan)

		lst := m.Mapping[info.RaceId]
		lst = m.deleteChanInList(sendChan, lst)

		if len(lst) > 0 {
			m.Mapping[info.RaceId] = lst
			println(`deleted player in race`, info.RaceId)
		} else {
			delete(m.Mapping, info.RaceId)
			println(`deleted race itself`, info.RaceId)
		}
	}
	return playerId
}
