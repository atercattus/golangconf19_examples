package scenes

import (
	"github.com/atercattus/golangconf19_examples/handmade/events"
	"github.com/atercattus/golangconf19_examples/handmade/net"
	"github.com/atercattus/golangconf19_examples/handmade/objects"
	"github.com/atercattus/golangconf19_examples/handmade/render"
	"github.com/atercattus/golangconf19_examples/handmade/resources"
	"math"
	"strconv"
	"syscall/js"
	"time"
)

type (
	SceneGame struct {
		SceneEmpty

		road *objects.Road

		batchRunBtn      *render.DrawCallBatch
		runBtn           *render.Sprite
		runBtnEtalonSize float32

		gophers            []*objects.Gopher
		gophersDrivingTime []float32

		gameStartedAt time.Time
		rhythm        float32
		rhythmSpeed   float32

		txtSpeedAndDistanceLastUpdate time.Time
		txtSpeedAndDistance           *render.Text2D

		alreadyFinished       bool
		resultPositionShown   bool
		gopherHasFinishedMark []bool
	}
)

var (
	_ Scener = &SceneGame{}
)

const (
	stepSpeedupInSteps  = 14
	stepSlowdownInSteps = 8

	rhythmMinSpeed = 1.0
	rhythmMaxSpeed = 3.0

	maxGopherSpeed = 100
)

func (scene *SceneGame) Show() {
	W, H := scene.renderer.Size()

	scene.txtSpeedAndDistance = scene.renderer.AddText(``, 40, render.Point{W / 2, 25})
	//scene.txtSpeedAndDistance.SetOutlineWidth(4)
	scene.updateSpeedAndDistanceTxt()

	players := make([]events.PlayerInfo, len(GlobalData.Players))
	for i := range GlobalData.Players {
		players[i] = GlobalData.Players[i].PlayerInfo
	}

	const roadY = 250
	scene.road = objects.NewRoad(scene.renderer, players, GlobalData.MaxPlayers, roadY)
	scene.road.SetFinishDistance(GlobalData.RaceDistance * GlobalData.GetPixelsPerStep())

	scene.gopherHasFinishedMark = make([]bool, len(GlobalData.Players))

	lineHeight := scene.road.GetLineHeightInPx()

	for idx := range GlobalData.Players {
		player := objects.NewGopher(scene.renderer)
		player.SizeTo(render.Point{Y: lineHeight * 0.6})
		player.MoveTo(render.Point{
			X: player.Size().X / 2,
			Y: scene.road.GetLineY(idx) - lineHeight*0.4,
		})
		scene.gophers = append(scene.gophers, player)
		scene.gophersDrivingTime = append(scene.gophersDrivingTime, 0)
	}

	scene.batchRunBtn = render.NewDrawCallBatch(scene.renderer, resources.TextureRunBtn)
	scene.runBtn = scene.batchRunBtn.AddSprite()
	scene.runBtnEtalonSize = scene.runBtn.Size().X
	scene.runBtn.MoveTo(render.Point{W / 2, H - scene.runBtnEtalonSize*1.1})

	GlobalData.WS.SetEventsReceiver(scene.serverEventsCb)

	scene.gameStartedAt = time.Now()
	scene.rhythm = 0
	scene.rhythmSpeed = rhythmMinSpeed
}

func (scene *SceneGame) Hide() {
	scene.renderer.DeleteText(scene.txtSpeedAndDistance)
}

func (scene *SceneGame) updateSpeedAndDistanceTxt() {
	playerInfo := &GlobalData.Players[GlobalData.PlayersOurIdx]
	toFinish := GlobalData.RaceDistance - playerInfo.Distance
	scene.txtSpeedAndDistance.SetText(
		`Скорость: ` + strconv.Itoa(int(playerInfo.Speed)) +
			`. До финиша:` + strconv.Itoa(int(toFinish)) + `м`,
	)
}

func (scene *SceneGame) serverEventsCb(event net.Event) {
	switch event.Command {
	case net.CommandRaceState:
		var ev events.EventRaceState
		if err := ev.Unmarshal([]byte(event.Data)); err != nil {
			println(`Wrong event:`, event, `Error:`, err.Error())
			return
		}
		ourFinishPos, smbdFinished := GlobalData.UpdatePlayersState(ev.Players)
		if ourFinishPos > 0 && scene.alreadyFinished && !scene.resultPositionShown {
			scene.resultPositionShown = true
			scene.txtSpeedAndDistance.SetText(`Вы финишировали ` + strconv.Itoa(ourFinishPos) + `м!`)
		}
		if smbdFinished {
			scene.finishedMarkProcess(ev)
		}

	case net.CommandEmpty:
		// pass

	default:
		println(`Unknown server event:`, event)
	}
}

func (scene *SceneGame) finishedMarkProcess(ev events.EventRaceState) {
	for i, playerState := range ev.Players {
		if playerState.FinishedAt > 0 && !scene.gopherHasFinishedMark[i] {
			scene.gopherHasFinishedMark[i] = true
			scene.road.SetPlayerName(i, `ФИНИШИРОВАЛ - `+GlobalData.Players[i].Name)
		}
	}
}

func (scene *SceneGame) OnClick(pos render.Point) {
	if !scene.alreadyFinished && scene.runBtn.IsPointInside(pos) {
		rhythmCoeff := scene.getCurrentRhythmCoeff()

		// При плохих тапах замедляемся, уходя в минус
		if thresh := float32(0.5); rhythmCoeff < thresh {
			rhythmCoeff -= 2 * thresh
		}

		speed := GlobalData.Players[GlobalData.PlayersOurIdx].Speed
		if speed += rhythmCoeff * stepSpeedupInSteps; speed < 0 {
			speed = 0
		}
		GlobalData.Players[GlobalData.PlayersOurIdx].Speed = speed

		if scene.rhythmSpeed += rhythmCoeff / 2.0; scene.rhythmSpeed > rhythmMaxSpeed {
			scene.rhythmSpeed = rhythmMaxSpeed
		} else if scene.rhythmSpeed < rhythmMinSpeed {
			scene.rhythmSpeed = rhythmMinSpeed
		}
	}
}

func (scene *SceneGame) updateGopherPos(idx int, posX float32, sprite *objects.Gopher) {
	//drivingTime := scene.gophersDrivingTime[idx]

	pos := sprite.Pos()
	pos.X = posX
	//pos.Y += float32(math.Sin(float64(drivingTime)) / 3.0) // небольшой сдвиг по вертикали для реализма
	sprite.MoveTo(pos)
}

func (scene *SceneGame) calcGopherSpeedAndDistance(dt float32, idx int, playerInfo *events.PlayerInfo) float32 {
	if playerInfo.Speed == 0 {
		return 0
	}

	scene.gophersDrivingTime[idx] += dt

	if playerInfo.Speed -= stepSlowdownInSteps * dt; playerInfo.Speed < 0 {
		playerInfo.Speed = 0
	} else if playerInfo.Speed > maxGopherSpeed {
		playerInfo.Speed = maxGopherSpeed
	}

	if playerInfo.Speed > 0 {
		delta := playerInfo.Speed * dt
		playerInfo.Distance += delta
		return delta
	}
	return 0
}

func (scene *SceneGame) updateOurGopherPos(dt float32) {
	idx := GlobalData.PlayersOurIdx
	playerInfo := &GlobalData.Players[idx]
	gopher := scene.gophers[idx]
	W, _ := scene.renderer.Size()

	pos := gopher.Pos()

	if playerInfo.Speed == 0 {
		// стоим на месте, движения нет
		scene.road.SetSpeed(0)
		return
	}

	if scene.isFinished(&playerInfo.PlayerInfo) {
		playerInfo.Speed = 0
		scene.rhythm = 0
		scene.rhythmSpeed = rhythmMinSpeed
		return
	}

	delta := scene.calcGopherSpeedAndDistance(dt, idx, &playerInfo.PlayerInfo)

	if pos.X < W/2 {
		// в самом начале пути, нужно доехать до середины экрана, при этом дорога еще не движется
		scene.road.SetSpeed(0)

		pos.X += delta * GlobalData.GetPixelsPerStep()
		scene.updateGopherPos(idx, pos.X, gopher)

		return
	}
	//// даже когда стоим на месте (катимся по середине дороги) применяю сдвиг по вертикали для красоты
	//scene.updateGopherPos(idx, pos.X, gopher)

	// спрайт больше не двигается, но двигается дорога
	scene.road.SetSpeed(playerInfo.Speed * GlobalData.GetPixelsPerStep())
}

func (scene *SceneGame) isFinished(playerInfo *events.PlayerInfo) bool {
	if scene.alreadyFinished {
		return true
	}

	if playerInfo.Distance >= GlobalData.RaceDistance {
		scene.txtSpeedAndDistance.SetText(`Вы финишировали!`)

		if err := GlobalData.SendIAmFinished(nil); err != nil {
			js.Global().Call(`alert`, `IAmFinished failed:`+err.Error())
		} else {
			scene.alreadyFinished = true
		}
		return true
	}
	return false
}

func (scene *SceneGame) updateOtherGopherPos(dt float32, idx int) {
	playerInfo := &GlobalData.Players[idx]

	delta := scene.calcGopherSpeedAndDistance(dt, idx, &playerInfo.PlayerInfo)
	distanceDiff := playerInfo.Distance - GlobalData.Players[GlobalData.PlayersOurIdx].Distance

	distanceDiff = (delta + distanceDiff) * GlobalData.GetPixelsPerStep()

	gopher := scene.gophers[idx]
	ourGopherPosX := scene.gophers[GlobalData.PlayersOurIdx].Pos().X
	scene.updateGopherPos(idx, ourGopherPosX+distanceDiff, gopher)
}

func (scene *SceneGame) Draw(dt float32) {
	if !scene.alreadyFinished {
		scene.rhythm += scene.rhythmSpeed * dt
		if scene.rhythmSpeed -= 0.2 * dt; scene.rhythmSpeed < rhythmMinSpeed {
			scene.rhythmSpeed = rhythmMinSpeed
		}
	}

	scene.road.Draw(dt)

	scene.updateOurGopherPos(dt)
	ourPlayerIdx := GlobalData.PlayersOurIdx
	for idx := range scene.gophers {
		if idx != ourPlayerIdx {
			scene.updateOtherGopherPos(dt, idx)
		}
	}

	for _, gopherSprite := range scene.gophers {
		gopherSprite.Draw()
	}

	if !scene.alreadyFinished {
		rhythmCoeff := scene.getCurrentRhythmCoeff()
		scene.runBtn.SetColor(1-rhythmCoeff, rhythmCoeff, 0, 1)
		scene.runBtn.SizeTo(render.Point{
			X: (rhythmCoeff + 1) * scene.runBtnEtalonSize,
		})
	}

	scene.batchRunBtn.Draw()

	if !scene.alreadyFinished {
		if now := time.Now(); scene.txtSpeedAndDistanceLastUpdate.Add(300 * time.Millisecond).Before(now) {
			scene.updateSpeedAndDistanceTxt()
			scene.txtSpeedAndDistanceLastUpdate = now
		}
	}
}

func (scene *SceneGame) getCurrentRhythmCoeff() float32 {
	intf, frac := math.Modf(float64(scene.rhythm))
	int_ := int(intf)
	if int_%2 == 1 {
		frac = 1 - frac
	}
	return float32(frac)
}
