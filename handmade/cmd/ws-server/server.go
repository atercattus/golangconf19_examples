package main

import (
	"context"
	"fmt"
	"github.com/atercattus/golangconf19_examples/handmade/events"
	"github.com/atercattus/golangconf19_examples/handmade/net"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type (
	CommandProcessor func(event net.Event, sendChan EventChan) net.Event

	WSServer struct {
		upgrader   websocket.Upgrader
		httpServer *http.Server

		commands map[net.Command]CommandProcessor

		gameState *GameState

		connections *RaceIdToPlayerConnections
	}
)

func MakeWSServer(addr string) *WSServer {
	srv := &WSServer{}
	srv.upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	srv.httpServer = &http.Server{Addr: addr, Handler: srv}

	srv.commands = map[net.Command]CommandProcessor{
		net.CommandPing: srv.processorPing,
		net.CommandJoin: srv.processorJoin,
		//net.CommandRaceState // Клиент это сам не запрашивает. Только сервер присылает.
		//net.CommandRaceInfo // Клиент это сам не запрашивает. Только сервер присылает.
		net.CommandIAmReady:    srv.processorIAmReady,
		net.CommandIAmFinished: srv.processorIAmFinished,
		net.CommandPlayerStep:  srv.processorPlayerStep,
	}

	srv.gameState = NewGameState(srv)

	srv.connections = NewRaceIdToPlayerConnections()

	return srv
}

func (s *WSServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(`WS: upgrade fail:`, err)
		return
	}
	defer conn.Close()

	clientUA := r.Header.Get(`User-Agent`)
	clientIP := r.Header.Get(`X-Real-Ip`)
	if clientIP == `` {
		clientIP = conn.RemoteAddr().String()
	}

	log.Printf(`New connect from %s (%s) / %s`, clientIP, conn.RemoteAddr().String(), clientUA)

	sendChan := make(EventChan, 100)

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		playerId := s.connections.DeleteConnectionBySendChan(sendChan)
		// если гонка еще не началась, то удаляю игрока из списка участников
		s.gameState.DeletePlayerFromFutureRace(playerId)

		cancel()
	}()

	go s.sender(ctx, cancel, conn, sendChan)

	helloEvent := net.NewEvent(
		net.CommandPing,
		fmt.Sprintf(`Hello %s / %s`, clientIP, clientUA),
	)
	sendChan <- &helloEvent

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if event, err := s.recvEvent(conn); err != nil {
			break
		} else {
			s.processEvent(conn, event, sendChan)
		}
	}
}

func (s *WSServer) sender(ctx context.Context, cancel context.CancelFunc, conn *websocket.Conn, ch EventChan) {
	for {
		select {
		case ev := <-ch:
			if err := s.sendEvent(conn, ev); err != nil {
				cancel()
				return
			}

		default:
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(50 * time.Millisecond)
			}
		}
	}
}

func (s *WSServer) recvEvent(conn *websocket.Conn) (event net.Event, err error) {
	if messageType, message, err := conn.ReadMessage(); err != nil {
		log.Printf("WS(%s): recv fail: %s\n", conn.RemoteAddr(), err)
		return event, err
	} else if messageType != websocket.TextMessage {
		return event, nil
	} else {
		err := event.Unmarshal(message)
		if err == nil {
			//log.Printf(`RECV %v`, event)
		}
		return event, err
	}
}

func (s *WSServer) sendEvent(conn *websocket.Conn, event *net.Event) error {
	err := conn.WriteMessage(websocket.TextMessage, event.Marshal())
	if err != nil {
		log.Printf("WS(%s): send fail: %s\n", conn.RemoteAddr(), err)
	}
	return err
}

func (s *WSServer) processEvent(conn *websocket.Conn, event net.Event, sendChan EventChan) {
	var resp net.Event

	if processor, ok := s.commands[event.Command]; ok {
		resp = processor(event, sendChan)
		resp.Idx = event.Idx // В ответе возвращаю номер из запроса
	} else {
		resp = net.NewCommandError(net.CommandErrorUnknownCommand)
	}

	select {
	case sendChan <- &resp:
	// pass
	case <-time.After(1 * time.Second):
		log.Printf(`Smth went wrong with sendChan when send %s to %s`, event.Marshal(), conn.RemoteAddr())
	}
}

func (s *WSServer) processorPing(event net.Event, sendChan EventChan) net.Event {
	event.Command = net.CommandPong
	return event
}

func (s *WSServer) processorJoin(event net.Event, sendChan EventChan) net.Event {
	var ev events.EventJoin
	if err := ev.Unmarshal([]byte(event.Data)); err != nil {
		return net.NewCommandError(err.Error())
	} else if race, err := s.gameState.JoinPlayer(ev); err != nil {
		return net.NewCommandError(err.Error())
	} else {
		s.connections.AddConnection(race.Id, ev.PlayerId, sendChan)
		return s.buildRaceNew(race)
	}
}

func (s *WSServer) processorIAmReady(event net.Event, sendChan EventChan) net.Event {
	var ev events.EventIAmReady
	if err := ev.Unmarshal([]byte(event.Data)); err != nil {
		return net.NewCommandError(err.Error())
	} else if err := s.gameState.IAmReady(ev); err != nil {
		return net.NewCommandError(err.Error())
	} else {
		return net.NewCommandEmpty()
	}
}

func (s *WSServer) processorIAmFinished(event net.Event, sendChan EventChan) net.Event {
	var ev events.EventIAmFinished
	if err := ev.Unmarshal([]byte(event.Data)); err != nil {
		return net.NewCommandError(err.Error())
	} else if err := s.gameState.IAmFinished(ev); err != nil {
		return net.NewCommandError(err.Error())
	} else {
		return net.NewCommandEmpty()
	}
}

func (s *WSServer) processorPlayerStep(event net.Event, sendChan EventChan) net.Event {
	var ev events.EventPlayerStep
	if err := ev.Unmarshal([]byte(event.Data)); err != nil {
		return net.NewCommandError(err.Error())
	} else if err := s.gameState.PlayerStep(ev); err != nil {
		return net.NewCommandError(err.Error())
	} else {
		return net.NewCommandEmpty()
	}
}

func (s *WSServer) buildRaceNew(race *Race) net.Event {
	resp := events.EventRaceNew{
		RaceUuid:   race.Id,
		MaxPlayers: maxPlayersInRace,
		Distance:   race.Distance,
	}
	return events.EventToNetEvent(&resp)
}

func (s *WSServer) buildRaceInfo(race *Race) net.Event {
	resp := events.EventRaceState{
		Players: race.GetPlayersState(),
	}
	return events.EventToNetEvent(&resp)
}

func (s *WSServer) buildRaceState(race *Race) net.Event {
	resp := events.EventRaceInfo{
		Players: race.GetPlayers(),
	}
	return events.EventToNetEvent(&resp)
}

func (s *WSServer) BroadcastEventToRace(raceId net.Id, ev net.Event) bool {
	lst := s.connections.GetChans(raceId)
	if len(lst) == 0 {
		return false
	}

	for _, sendChan := range lst {
		sendChan <- &ev
	}

	return true
}
