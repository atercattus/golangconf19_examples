package scenes

import (
	"github.com/EngoEngine/engo"
	engoCommon "github.com/EngoEngine/engo/common"
	"github.com/atercattus/golangconf19_examples/engo/common"
	"github.com/atercattus/golangconf19_examples/engo/systems"
)

type (
	TitleScene struct {
		common.CommonScene
		gopher *systems.GopherSystem
	}
)

const (
	imgFilenameGoRaces = `images/go_races.png`
)

func (scene *TitleScene) Preload() {
	scene.CommonScene.Preload()
	_ = engo.Files.Load(imgFilenameGoRaces)
}

func (scene *TitleScene) Setup(u engo.Updater) {
	scene.CommonScene.Setup(u)

	//engo.SetScene(Scenes.SelectName, false) // ToDo: дебаг

	if textureGopher, err := engoCommon.LoadedSprite(imgFilenameGoRaces); err != nil {
		panic(err)
	} else {
		common.NewSprite(textureGopher, &common.SpriteOptions{
			Position: engo.Point{
				X: (engo.WindowWidth() - textureGopher.Width()) / 2,
				Y: engo.GameHeight()/2 - textureGopher.Height(),
			},
		}).AddToWorld(scene.World())
	}

	gopher := systems.NewGopherSystem(scene.World())
	scene.gopher = gopher

	gopher.MoveTo(engo.Point{
		X: -scene.gopher.Width(),
		Y: engo.WindowHeight() * 5 / 8,
	})
	gopher.SetRelativeSpeed(engo.WindowWidth() / 2)
	gopher.SetIsDriving(true)
	gopher.SetUpdateCb(scene.GopherUpdateCb)

	// about me :)
	common.NewText(`created by AterCattus`, common.Font80, &common.TextOptions{
		SpriteOptions: common.SpriteOptions{
			Position: engo.Point{
				//X: 0,
				Y: engo.WindowHeight() - float32(common.Font80Size)*0.25 - 10,
			},
			Scale:  0.25,
			Width:  engo.WindowWidth(),
			Height: engo.WindowHeight(),
		},
		AlignH: common.TextAlignMin,
	}).AddToWorld(u)
}

func (scene *TitleScene) GopherUpdateCb(dt float32, gopher *systems.GopherSystem) {
	if gopher.Position().X > engo.WindowWidth() {
		engo.SetScene(Scenes.SelectName, false)
	}
}
