package scenes

import (
	"github.com/atercattus/golangconf19_examples/handmade/render"
)

type (
	SceneEmpty struct {
		renderer     *render.WebGLRender
		sceneManager *SceneManager
	}
)

var (
	_ Scener = &SceneEmpty{}
)

func (scene *SceneEmpty) Setup(renderer *render.WebGLRender, sceneManager *SceneManager) {
	scene.renderer = renderer
	scene.sceneManager = sceneManager
}

func (scene *SceneEmpty) Show() {
}

func (scene *SceneEmpty) Hide() {
}

func (scene *SceneEmpty) Draw(dt float32) {
}

func (scene *SceneEmpty) OnClick(pos render.Point) {
}
