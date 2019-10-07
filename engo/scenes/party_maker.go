package scenes

import (
	"github.com/EngoEngine/engo"
	"github.com/atercattus/golangconf19_examples/engo/common"
	"github.com/atercattus/golangconf19_examples/engo/systems"
	"time"
)

type (
	PartyMakerScene struct {
		common.CommonScene
	}
)

func (scene *PartyMakerScene) Preload() {
	scene.CommonScene.Preload()
}

func (scene *PartyMakerScene) Setup(u engo.Updater) {
	scene.CommonScene.Setup(u)

	systems.After(scene.World(), 1*time.Second, func() {
		engo.SetScene(Scenes.Game, false)
	})
}
