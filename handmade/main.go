package main

import (
	"context"
	"github.com/atercattus/golangconf19_examples/handmade/net"
	"github.com/atercattus/golangconf19_examples/handmade/render"
	"github.com/atercattus/golangconf19_examples/handmade/scenes"
	"github.com/nuberu/webgl"
	"log"
	"math/rand"
	"syscall/js"
	"time"
)

var (
	gl       *webgl.RenderingContext
	renderer *render.WebGLRender
	sceneMan *scenes.SceneManager
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go _main(ctx) // иначе можно словить "fatal error: all goroutines are asleep - deadlock!"

	<-make(chan struct{})
}

func _main(ctx context.Context) {
	rand.Seed(time.Now().UnixNano())

	scenes.GlobalData.Init()

	go func() {
		scenes.GlobalData.WS = net.NewWebsocketClient()

		go pinger(ctx)

		const wsHost = `wss://ater.me/go_races/ws`
		//const wsHost = `ws://127.0.0.1:8081/`
		if err := scenes.GlobalData.WS.DialAndServe(ctx, wsHost); err != nil {
			log.Println(`WS error:`, err)
		}
	}()

	var err error
	renderer, err = render.NewWebGLRender(`canvas`, `canvas_text`)
	if err != nil {
		js.Global().Call(`alert`, err) // console.log ?
		return
	}

	W, H := renderer.Size()

	scenes.GlobalData.SetPixelsPerStep(W / scenes.ScreenWidthInSteps)

	gl = renderer.GetGlCtx()

	gl.ClearColor(1, 1, 1, 1)
	gl.Clear(uint32(webgl.COLOR_BUFFER_BIT))

	gl.Enable(webgl.BLEND)
	gl.BlendFunc(webgl.SRC_ALPHA, webgl.ONE_MINUS_SRC_ALPHA)

	gl.Viewport(0, 0, int(W), int(H))

	sceneMan = scenes.NewSceneManager(renderer)
	sceneMan.Register(`title`, &scenes.SceneTitle{})
	sceneMan.Register(`select_name`, &scenes.SceneSelectName{})
	sceneMan.Register(`party_maker`, &scenes.ScenePartyMaker{})
	sceneMan.Register(`game`, &scenes.SceneGame{})
	sceneMan.Goto(`title`)

	renderer.SetRenderFunc(func(dt float32) {
		renderer.Clear()
		sceneMan.Current().Draw(dt)
	})

	renderer.SetClickFunc(func(pos render.Point) {
		sceneMan.Current().OnClick(pos)
	})
}

func pinger(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.Tick(10 * time.Second):
			pingEvent := net.NewEvent(net.CommandPing, string(scenes.GlobalData.PlayerId))
			_ = scenes.GlobalData.WS.SendMessage(pingEvent, func(net.Event) {})
		}
	}
}
