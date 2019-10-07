package scenes

import (
	"github.com/atercattus/golangconf19_examples/handmade/events"
	"github.com/atercattus/golangconf19_examples/handmade/net"
	"github.com/atercattus/golangconf19_examples/handmade/render"
	"github.com/atercattus/golangconf19_examples/handmade/resources"
	"math"
	"syscall/js"
)

type (
	charBBox struct {
		X1, Y1, X2, Y2 float32
	}
	charInfo struct {
		char string
		text *render.Text2D
		bbox charBBox
	}

	SceneSelectName struct {
		SceneEmpty

		titleText *render.Text2D
		chars     []charInfo

		playerName string

		batchNextBtn *render.DrawCallBatch
		nextBtn      *render.Sprite

		batchBtn *render.DrawCallBatch
	}
)

var (
	_ Scener = &SceneSelectName{}
)

const (
	titleTextPrefix  = `Your name: `
	backspaceChar    = `<<`
	playernameMaxLen = 10
)

func (scene *SceneSelectName) nextBtnShow(show bool) {
	W, H := scene.renderer.Size()

	// show/hide пока не работают. Так что просто двигаю кнопку по вертикали за пределы экрана.

	if show {
		scene.nextBtn.MoveTo(render.Point{W / 2, H - scene.nextBtn.Size().Y})
	} else {
		scene.nextBtn.MoveTo(render.Point{W / 2, 2 * H})
	}
}

func (scene *SceneSelectName) Show() {
	W, H := scene.renderer.Size()

	fontSize := int(math.Min(float64(W)/15, float64(H)/15))

	scene.titleText = scene.renderer.AddText(titleTextPrefix, fontSize,
		render.Point{
			X: float32(fontSize) / 2,
			Y: float32(fontSize),
		},
	)
	scene.titleText.SetAlignX(`left`)

	// next
	scene.batchNextBtn = render.NewDrawCallBatch(scene.renderer, resources.TextureRunBtn)
	scene.nextBtn = scene.batchNextBtn.AddSprite()
	scene.nextBtn.SizeTo(render.Point{X: scene.nextBtn.Size().X * 1.5})
	nextBtnY := H - scene.nextBtn.Size().Y
	scene.nextBtnShow(false)

	scene.batchBtn = render.NewDrawCallBatch(scene.renderer, resources.TextureBtn)

	charsTopLineY := scene.titleText.Pos().Y + 2*scene.titleText.Size()

	const charsCount = 'Z' - 'A' + 1
	const charsInLine = 6
	const charsLines = (charsCount + (charsInLine - 1)) / charsInLine
	charW := W / charsInLine
	charH := (nextBtnY - charsTopLineY) / charsLines

	// chars
	pos := render.Point{
		X: float32(fontSize),
		Y: charsTopLineY,
	}
	scene.chars = make([]charInfo, 0, charsCount+1) // +1 для backspace
	for ch := 'A'; ch <= 'Z'; ch++ {
		char := scene.renderer.AddText(string(ch), fontSize, pos)
		scene.chars = append(scene.chars, charInfo{
			char: string(ch),
			text: char,
			bbox: charBBox{
				X1: pos.X - charW/2,
				Y1: pos.Y - charH/2,
				X2: pos.X + charW/2,
				Y2: pos.Y + charH/2,
			},
		})
		sprite := scene.batchBtn.AddSprite()
		sprite.MoveTo(pos)
		sprite.SizeTo(render.Point{X: sprite.Size().X * 1.2})

		if pos.X += charW; pos.X >= W {
			pos.X = float32(fontSize)
			pos.Y += charH
		}
	}

	// backspace
	pos.X = float32(fontSize) + charW*(charsInLine-1)
	scene.chars = append(scene.chars, charInfo{
		char: backspaceChar,
		text: scene.renderer.AddText(backspaceChar, fontSize, pos),
		bbox: charBBox{
			X1: pos.X - charW/2,
			Y1: pos.Y - charH/2,
			X2: pos.X + charW/2,
			Y2: pos.Y + charH/2,
		},
	})
	sprite := scene.batchBtn.AddSprite()
	sprite.MoveTo(pos)
	sprite.SizeTo(render.Point{X: sprite.Size().X * 1.2})
}

func (scene *SceneSelectName) Hide() {
	scene.renderer.DeleteText(scene.titleText)
	for _, charInfo := range scene.chars {
		scene.renderer.DeleteText(charInfo.text)
	}
	scene.chars = nil
}

func (scene *SceneSelectName) OnClick(pos render.Point) {
	if scene.nextBtn.IsPointInside(pos) {
		scene.join()
		return
	}

	for _, charInfo := range scene.chars {
		bbox := charInfo.bbox
		if pos.X >= bbox.X1 && pos.X <= bbox.X2 && pos.Y >= bbox.Y1 && pos.Y <= bbox.Y2 {
			scene.updatePlayerName(charInfo.char)
			break
		}
	}
}

func (scene *SceneSelectName) updatePlayerName(newChar string) {
	if newChar == backspaceChar {
		if len(scene.playerName) == 0 {
			return
		}
		scene.playerName = scene.playerName[:len(scene.playerName)-1]
	} else {
		if len(scene.playerName) >= playernameMaxLen {
			return
		}
		scene.playerName = scene.playerName + newChar
	}

	scene.titleText.SetText(titleTextPrefix + scene.playerName)

	if len(scene.playerName) < 3 {
		scene.nextBtnShow(false)
	} else {
		scene.nextBtnShow(true)
	}
}

func (scene *SceneSelectName) Draw(dt float32) {
	scene.batchBtn.Draw()

	scene.batchNextBtn.Draw()
}

func (scene *SceneSelectName) join() {
	cb := func(event net.Event) {
		var ev events.EventRaceNew
		if err := ev.Unmarshal([]byte(event.Data)); err != nil {
			js.Global().Call(`alert`, `Join got wrong resp:`+err.Error())
			return
		}
		GlobalData.RaceId = ev.RaceUuid
		GlobalData.RaceDistance = ev.Distance
		GlobalData.MaxPlayers = ev.MaxPlayers

		scene.sceneManager.Goto(`party_maker`)
	}

	if err := GlobalData.Join(scene.playerName, cb); err != nil {
		js.Global().Call(`alert`, `Join failed:`+err.Error())
	}
}
