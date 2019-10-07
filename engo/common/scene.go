package common

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"image/color"
)

type (
	SceneName string

	CommonScene struct {
		engoUpdater engo.Updater
		Name        SceneName
	}
)

var (
	_ engo.Scene = &CommonScene{}
)

func (scene *CommonScene) Type() string {
	return string(scene.Name)
}

func (*CommonScene) Preload() {
	FontsPreload()
}

func (scene *CommonScene) Setup(u engo.Updater) {
	scene.engoUpdater = u

	common.SetBackground(color.White)

	world, _ := u.(*ecs.World)
	world.AddSystem(&common.RenderSystem{})
	//world.AddSystem(&common.FPSSystem{Terminal: true})
}

func (scene *CommonScene) World() *ecs.World {
	world, _ := scene.engoUpdater.(*ecs.World)
	return world
}
