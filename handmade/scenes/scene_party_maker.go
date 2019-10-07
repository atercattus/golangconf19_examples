package scenes

import (
	"github.com/atercattus/golangconf19_examples/handmade/events"
	"github.com/atercattus/golangconf19_examples/handmade/net"
	"github.com/atercattus/golangconf19_examples/handmade/objects"
	"github.com/atercattus/golangconf19_examples/handmade/render"
	"github.com/atercattus/golangconf19_examples/handmade/resources"
	"strconv"
	"syscall/js"
)

type (
	ScenePartyMaker struct {
		SceneEmpty

		road *objects.Road

		batchRunBtn *render.DrawCallBatch
		runBtn      *render.Sprite

		txtTitle    *render.Text2D
		txtDistance *render.Text2D
		txtRunDescr *render.Text2D
	}
)

var (
	_ Scener = &ScenePartyMaker{}
)

func (scene *ScenePartyMaker) Show() {
	W, H := scene.renderer.Size()

	scene.txtTitle = scene.renderer.AddText(`Сбор участников`, 100, render.Point{W / 2, 60})

	distStr := `Длина трассы: ` + strconv.Itoa(int(GlobalData.RaceDistance)) + ` м`
	scene.txtDistance = scene.renderer.AddText(distStr, 60, render.Point{W / 2, 140})

	scene.road = objects.NewRoad(scene.renderer, nil, GlobalData.MaxPlayers, 250)
	scene.road.SetSpeed(GlobalData.GetPixelsPerStep() * ScreenWidthInSteps / 20.0) // ширину экрана за 20 секунд

	scene.batchRunBtn = render.NewDrawCallBatch(scene.renderer, resources.TextureRunBtn)
	scene.runBtn = scene.batchRunBtn.AddSprite()
	scene.runBtn.SizeTo(render.Point{X: scene.runBtn.Size().X * 2})
	runBtnY := H - scene.runBtn.Size().Y*0.7
	scene.runBtn.MoveTo(render.Point{W / 2, runBtnY})

	const txtSize = 50
	txtY := scene.runBtn.Pos().Y - scene.runBtn.Size().Y/2 - txtSize/2
	scene.txtRunDescr = scene.renderer.AddText(`Жми по готовности:`, txtSize, render.Point{W / 2, txtY})

	GlobalData.WS.SetEventsReceiver(scene.serverEventsCb)
}

func (scene *ScenePartyMaker) Hide() {
	GlobalData.WS.SetEventsReceiver(nil)

	scene.road.Delete(scene.renderer)

	scene.renderer.DeleteText(scene.txtTitle)
	scene.renderer.DeleteText(scene.txtDistance)
	scene.renderer.DeleteText(scene.txtRunDescr)
}

func (scene *ScenePartyMaker) OnClick(pos render.Point) {
	if scene.runBtn.IsPointInside(pos) {
		cb := func(event net.Event) {
			scene.runBtn.SetColor(0.2, 1, 0.2, 1)
			scene.txtRunDescr.SetText(`Ждем остальных...`)
		}

		scene.runBtn.SetColor(0.8, 0.8, 0.8, 0.8)
		scene.txtRunDescr.SetText(`Запрос к серверу...`)

		if err := GlobalData.SendIAmReady(cb); err != nil {
			js.Global().Call(`alert`, `IAmReady failed:`+err.Error())
		}
		return
	}
}

func (scene *ScenePartyMaker) serverEventsCb(event net.Event) {
	switch event.Command {
	case net.CommandRaceInfo:
		var ev events.EventRaceInfo
		if err := ev.Unmarshal([]byte(event.Data)); err != nil {
			println(`Wrong event:`, event, `Error:`, err.Error())
			return
		}
		GlobalData.SetPlayersList(ev.Players)
		scene.updateRaceInfo()

	case net.CommandRaceStart:
		var ev events.EventRaceStart
		if err := ev.Unmarshal([]byte(event.Data)); err != nil {
			println(`Wrong event:`, event, `Error:`, err.Error())
			return
		}
		GlobalData.RaceId = ev.RaceId
		scene.sceneManager.Goto(`game`)

	default:
		println(`Unknown server event:`, event)
	}
}

func (scene *ScenePartyMaker) updateRaceInfo() {
	for i, playerInfo := range GlobalData.Players {
		name := playerInfo.Name
		if playerInfo.Ready {
			name += ` - В ДЕЛЕ` // ToDo:
		} else {
			name += ` - ОЖИДАЕМ`
		}
		scene.road.SetPlayerName(i, name)
	}
	for i := len(GlobalData.Players); i < GlobalData.MaxPlayers; i++ {
		scene.road.SetPlayerName(i, ``)
	}
}

func (scene *ScenePartyMaker) Draw(dt float32) {
	if scene.road != nil {
		scene.road.Draw(dt)
	}
	scene.batchRunBtn.Draw()
}
