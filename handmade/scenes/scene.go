package scenes

import "github.com/atercattus/golangconf19_examples/handmade/render"

type (
	Scener interface {
		Setup(renderer *render.WebGLRender, sceneManager *SceneManager)
		Show()
		Hide()
		Draw(dt float32)
		OnClick(pos render.Point)
	}
)
