// +build js

package net

import (
	"context"
	"log"
	"syscall/js"
)

type (
	WebsocketClient struct {
		ws          js.Value
		wsConnected bool
		urlStr      string

		callbacks map[EventId]EventCallback
	}
)

var (
	_ WebsocketClienter = &WebsocketClient{}
)

func NewWebsocketClient() *WebsocketClient {
	return &WebsocketClient{
		callbacks: map[EventId]EventCallback{},
	}
}

func (ws *WebsocketClient) IsConnected() bool {
	return ws.wsConnected
}

func (ws *WebsocketClient) DialAndServe(ctx context.Context, urlStr string) error {
	ws.urlStr = urlStr
	ws.reconnect()
	defer ws.closeConn()

	<-ctx.Done()

	return nil
}

func (ws *WebsocketClient) SendMessage(event Event, cb EventCallback) (err error) {
	if !ws.wsConnected {
		ws.reconnect()
		return ErrWSIsNotConnected
	}

	defer func() {
		if e := recover(); e == nil {
		} else if jsErr, ok := e.(*js.Error); ok && jsErr != nil {
			err = jsErr
		} else if stdErr, ok := e.(error); ok && stdErr != nil {
			err = stdErr
		} else {
			log.Println("WS send error:", e)
		}
	}()

	ws.ws.Call(`send`, string(event.Marshal()))

	// Это выполняется в однопоточном js,
	//   и тут мы оказываемся до того, как от сервера может придти ответ.
	if cb != nil {
		ws.callbacks[event.Idx] = cb
	}

	return
}

func (ws *WebsocketClient) closeConn() {
	if ws.wsConnected {
		ws.ws.Call(`close`, 1000, `Closing Connection Normally`)
		ws.wsConnected = false
	}
}

func (ws *WebsocketClient) reconnect() (success bool) {
	success = false
	if ws.urlStr == `` {
		return
	}

	defer func() {
		if e := recover(); e != nil {
			log.Println("WS connect error:", e)
		}
	}()

	ws.ws = js.Global().Get(`WebSocket`).New(ws.urlStr)

	ws.ws.Call(`addEventListener`, `open`, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ws.wsConnected = true
		return nil
	}))

	ws.ws.Call(`addEventListener`, `message`, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		data := e.Get(`data`).String()
		event := NewEvent(CommandEmpty, ``)
		if err := event.Unmarshal([]byte(data)); err != nil {
			log.Printf(`WS wrong event(%s): %s`, data, err)
		} else if asyncCb, ok := ws.callbacks[event.Idx]; ok {
			delete(ws.callbacks, event.Idx)
			asyncCb(event)
		} // иначе непонятный колбек вникуда (уже обработали?)
		return nil
	}))

	ws.ws.Call(`addEventListener`, `error`, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		log.Println(`WS error:`, e.Get(`message`).String())
		return nil
	}))

	ws.ws.Call(`addEventListener`, `close`, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		log.Printf(`WS closed. clean:%v code:%d reason:%s`,
			e.Get(`wasClean`).Bool(),
			e.Get(`code`).Int(),
			e.Get(`reason`).String(),
		)
		ws.closeConn()
		return nil
	}))

	return true
}
