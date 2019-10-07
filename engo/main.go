package main

import (
	"context"
	"github.com/EngoEngine/engo"
	"github.com/atercattus/golangconf19_examples/engo/common"
	"github.com/atercattus/golangconf19_examples/engo/net"
	"github.com/atercattus/golangconf19_examples/engo/scenes"
	"log"
	"math/rand"
	"runtime"
	"time"
)

const (
	screenWidth  = 600 //576
	screenHeight = 800 //1024
)

func main() {
	rand.Seed(time.Now().UnixNano())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		common.GlobalData.WS = net.NewWebsocketClient()
		if err := common.GlobalData.WS.DialAndServe(ctx, `ws://127.0.0.1:8081/`); err != nil {
			log.Println(`WS error:`, err)
		}
	}()

	//go func() {
	//	for range time.Tick(1 * time.Second) {
	//		event := net.NewEvent(net.CommandPing, ``)
	//		event.Data = time.Now().Format(time.RFC3339)
	//		err := common.GlobalData.WS.SendMessage(event, func(event net.Event) {
	//			log.Println("WS ping response:", event)
	//		})
	//		if err != nil {
	//			log.Println(`WS send error:`, err)
	//		}
	//	}
	//}()

	// в нативном приложении тут будут нули, а в gopherjs - размеры области
	w := int(engo.WindowWidth())
	if w == 0 {
		w = screenWidth
	}
	h := int(engo.WindowHeight())
	if h == 0 {
		h = screenHeight
	}

	// для загрузки графики и нативно и в вебе
	assetsRoot := `assets`
	if runtime.GOARCH == `js` || runtime.GOARCH == `wasm` {
		assetsRoot = `.`
	}

	opts := engo.RunOptions{
		Title:        `GolangConf Engo demo`,
		Width:        w,
		Height:       h,
		NotResizable: true,
		VSync:        true,
		AssetsRoot:   assetsRoot,
		//Update: // !!!
	}
	engo.Run(opts, scenes.Scenes.Title)
}
