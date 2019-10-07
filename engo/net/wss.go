// +build !js

package net

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

type (
	WebsocketClient struct {
		byteSlicePool sync.Pool
		sendChan      chan []byte
		urlStr        string
		wsConn        *websocket.Conn

		callbacks map[EventId]EventCallback
	}
)

var (
	_ WebsocketClienter = &WebsocketClient{}

	ErrSendChanOverflow = fmt.Errorf(`WS send chan overflow`)
)

func NewWebsocketClient() (ws *WebsocketClient) {
	ws = &WebsocketClient{}

	ws.sendChan = make(chan []byte, 1000)

	ws.byteSlicePool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 1024)
		},
	}

	ws.callbacks = map[EventId]EventCallback{}

	return ws
}

func (ws *WebsocketClient) SendMessage(event Event, cb EventCallback) error {
	tmp := ws.byteSlicePool.Get().([]byte)
	tmp = append(tmp[:0], event.Marshal()...)

	if cb != nil {
		ws.callbacks[event.Idx] = cb
	}

	select {
	case ws.sendChan <- tmp:
		return nil
	default:
		ws.byteSlicePool.Put(tmp)
		delete(ws.callbacks, event.Idx)
		return ErrSendChanOverflow
	}
}

func (ws *WebsocketClient) IsConnected() bool {
	return ws.wsConn != nil
}

func (ws *WebsocketClient) DialAndServe(ctx context.Context, urlStr string) error {
	ws.urlStr = urlStr
	ws.reconnect()
	defer ws.closeConn()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			event := NewEvent(CommandEmpty, ``)

			if wsConn := ws.wsConn; wsConn == nil {
				time.Sleep(100 * time.Millisecond)
			} else if _, message, err := wsConn.ReadMessage(); err != nil {
				log.Println(`WS read error:`, err)
				time.Sleep(100 * time.Millisecond)
			} else if err := event.Unmarshal(message); err != nil {
				log.Printf(`WS wrong event(%s): %s`, message, err)
			} else if asyncCb, ok := ws.callbacks[event.Idx]; ok {
				delete(ws.callbacks, event.Idx)
				asyncCb(event)
			} // иначе непонятный колбек вникуда (уже обработали?)
		}
	}()

	sendLoopTicker := time.Tick(1 * time.Second)
sendLoop:
	for {
		if (ws.wsConn == nil) && !ws.reconnect() {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		select {
		case data := <-ws.sendChan:
			if err := ws.wsConn.WriteMessage(websocket.TextMessage, data); err == nil {
			} else if !ws.reconnect() {
				log.Println(`WS send error:`, err)
			} else if err := ws.wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println(`WS send after reconnect error:`, err)
			}
			ws.byteSlicePool.Put(data)

		case <-ctx.Done():
			break sendLoop

		case <-sendLoopTicker:
			// pass
		}
	}

	wg.Wait()

	return nil
}

func (ws *WebsocketClient) closeConn() {
	if ws.wsConn != nil {
		_ = ws.wsConn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ``),
		)
		_ = ws.wsConn.Close()
		ws.wsConn = nil
	}
}

func (ws *WebsocketClient) reconnect() bool {
	ws.closeConn()

	if ws.urlStr == `` {
		return false
	}

	dialCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	conn, _, err := websocket.DefaultDialer.DialContext(dialCtx, ws.urlStr, nil)
	if err != nil {
		log.Println(`WS connect error:`, err)
		return false
	}
	ws.wsConn = conn

	return true
}
