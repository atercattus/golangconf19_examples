package scenes

import (
	"github.com/EngoEngine/engo"
	engoCommon "github.com/EngoEngine/engo/common"
	"github.com/atercattus/golangconf19_examples/engo/common"
	"github.com/atercattus/golangconf19_examples/engo/events"
	"github.com/atercattus/golangconf19_examples/engo/net"
	"github.com/atercattus/golangconf19_examples/engo/systems"
	"log"
	"time"
)

type (
	SelectNameScene struct {
		common.CommonScene
		nameText *common.Text
		applyBtn *common.Text
	}
)

var (
	backspaceStr   = `<<`
	titleHeight    = float32(common.Font80Size) * 0.9
	nameTextHeight = float32(common.Font80Size) * 0.9
	charsTopLineY  = (titleHeight + nameTextHeight) * 1.2
)

func (scene *SelectNameScene) Preload() {
	scene.CommonScene.Preload()
}

func (scene *SelectNameScene) Setup(u engo.Updater) {
	scene.CommonScene.Setup(u)

	mouseSystem := &engoCommon.MouseSystem{}
	scene.World().AddSystem(mouseSystem)

	w := engo.WindowWidth()
	const fontScale = 0.7
	const charsInLine = 6
	charW := w / charsInLine

	// title
	common.NewText(`Your name:`, common.Font80, &common.TextOptions{
		SpriteOptions: common.SpriteOptions{
			Width:  w,
			Height: titleHeight,
			Scale:  0.7,
		},
		AlignH: common.TextAlignCenter,
		AlignV: common.TextAlignCenter,
	}).AddToWorld(u)

	// player name
	scene.nameText = common.NewText(``, common.Font80DarkBlue, &common.TextOptions{
		SpriteOptions: common.SpriteOptions{
			Position: engo.Point{
				X: 0,
				Y: titleHeight,
			},
			Width:  engo.WindowWidth(),
			Height: nameTextHeight,
			Scale:  0.9,
		},
		AlignH: common.TextAlignCenter,
	})
	scene.nameText.AddToWorld(u)

	// chars
	var pos engo.Point
	pos.Y = charsTopLineY
	for ch := 'A'; ch <= 'Z'; ch++ {
		char := common.NewText(string(ch), common.Font80, &common.TextOptions{
			SpriteOptions: common.SpriteOptions{
				Position: pos,
				Width:    charW,
				Height:   float32(common.Font80Size),
				Scale:    float32(fontScale),
			},
			AlignH: common.TextAlignCenter,
			AlignV: common.TextAlignCenter,
		})
		char.AddToWorld(u)

		mouseSystem.Add(&char.BasicEntity, &char.MouseComponent, &char.SpaceComponent, &char.RenderComponent)
		scene.World().AddSystem(&systems.MousableSystem{Callback: func(dt float32) {
			if char.MouseComponent.Clicked {
				scene.charCallback(char, dt)
			}
		}})

		if pos.X += charW; pos.X >= w {
			pos.X = 0
			pos.Y += float32(common.Font80Size)
		}
	}

	// backspace
	if true {
		char := common.NewText(backspaceStr, common.Font80, &common.TextOptions{
			SpriteOptions: common.SpriteOptions{
				Position: engo.Point{
					X: charW * (charsInLine - 1),
					Y: pos.Y,
				},
				Width:  charW,
				Height: float32(common.Font80Size),
				Scale:  float32(fontScale),
			},
			AlignH: common.TextAlignMin,
		})
		char.AddToWorld(u)
		mouseSystem.Add(&char.BasicEntity, &char.MouseComponent, &char.SpaceComponent, &char.RenderComponent)
		scene.World().AddSystem(&systems.MousableSystem{Callback: func(dt float32) {
			if char.MouseComponent.Clicked {
				scene.charCallback(char, dt)
			}
		}})
	}

	// apply
	scene.applyBtn = common.NewText(`Next`, common.Font80, &common.TextOptions{
		SpriteOptions: common.SpriteOptions{
			Position: engo.Point{
				X: 0,
				Y: pos.Y + float32(common.Font80Size)*1.3,
			},
			Width:  engo.WindowWidth(),
			Height: nameTextHeight,
			Scale:  0.9,
		},
		AlignH: common.TextAlignCenter,
	})
	scene.applyBtn.AddToWorld(u)
	scene.applyBtn.Hidden = true
	mouseSystem.Add(&scene.applyBtn.BasicEntity, &scene.applyBtn.MouseComponent, &scene.applyBtn.SpaceComponent, &scene.applyBtn.RenderComponent)
	scene.World().AddSystem(&systems.MousableSystem{Callback: func(dt float32) {
		if scene.applyBtn.MouseComponent.Clicked {
			scene.join()
		}
	}})
}

func (scene *SelectNameScene) charCallback(char *common.Text, dt float32) {
	charText := char.Drawable.(engoCommon.Text).Text

	txt := scene.nameText.Drawable.(engoCommon.Text)
	if charText == backspaceStr {
		if len(txt.Text) > 0 {
			txt.Text = txt.Text[:len(txt.Text)-1]
		} else {
			return
		}
	} else if len(txt.Text) < 10 {
		txt.Text += charText
	} else {
		return
	}
	scene.nameText.Drawable = txt

	scene.applyBtn.Hidden = len(txt.Text) < 3

	scale := float32(0.9)
	if len(txt.Text) > 8 {
		scale = 0.75
	}
	scene.nameText.Scale.X = scale
	scene.nameText.Scale.Y = scale

	// Realign() не знает исходного положения текста, так что возвращаю его на место
	scene.nameText.Position.Y = titleHeight
	scene.nameText.Position.X = 0
	scene.nameText.Realign()
}

func (scene *SelectNameScene) join() {
	username := scene.nameText.Drawable.(engoCommon.Text).Text

	cb := func(event net.Event) {
		var ev events.EventUsersCount
		if err := ev.Unmarshal([]byte(event.Data)); err != nil {
			log.Println("Join got wrong resp:", err)
			time.Sleep(100 * time.Millisecond)
			scene.join()
			return
		}
		common.GlobalData.UsersCount = ev.Count

		engo.SetScene(Scenes.PartyMaker, false)
	}

	if err := common.GlobalData.Join(username, cb); err != nil {
		log.Println("Join fail:", err)
		time.Sleep(100 * time.Millisecond)
		scene.join()
	}
}
