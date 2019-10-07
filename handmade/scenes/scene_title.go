package scenes

import (
	"github.com/atercattus/golangconf19_examples/handmade/objects"
	"github.com/atercattus/golangconf19_examples/handmade/render"
	"github.com/atercattus/golangconf19_examples/handmade/resources"
	"math"
)

type (
	SceneTitle struct {
		SceneEmpty
		sceneShown bool

		iter int

		batchTitle *render.DrawCallBatch

		gopher    *objects.Gopher
		driveTime float32

		loadingText *render.Text2D
	}
)

var (
	_ Scener = &SceneTitle{}
)

func (scene *SceneTitle) Show() {
	W, H := scene.renderer.Size()

	scene.loadingText = scene.renderer.AddText(`ver 2`, 40, render.Point{W * 0.5, H - 50})
	scene.loadingText.SetOutlineWidth(10)

	resources.TexturesPreload(scene.renderer, func() {
		scene.batchTitle = render.NewDrawCallBatch(scene.renderer, resources.TextureTitle)
		if sprite := scene.batchTitle.AddSprite(); true {
			sprite.SizeTo(render.Point{W * 0.8, 0})
			sprite.MoveTo(render.Point{W * 0.5, H * 0.3})
		}

		scene.gopher = objects.NewGopher(scene.renderer)

		scene.gopher.SizeTo(render.Point{W * 0.2, 0})
		scene.gopher.MoveTo(render.Point{-scene.gopher.Size().X / 2, H * 0.6})

		scene.sceneShown = true
	})
}

func (scene *SceneTitle) Hide() {
	scene.renderer.DeleteText(scene.loadingText)
}

func (scene *SceneTitle) Draw(dt float32) {
	if !scene.sceneShown {
		return
	}

	speed := scene.renderer.Width() / 2
	scene.driveTime += dt

	scene.gopher.MoveBy(render.Point{
		X: speed * dt,
		Y: float32(math.Sin(float64(scene.driveTime)) / 2.0), // небольшой сдвиг по вертикали для реализма
	})
	if scene.gopher.Pos().X > scene.renderer.Width() {
		scene.sceneManager.Goto(`select_name`)
		return
	}

	scene.batchTitle.Draw()
	scene.gopher.Draw()
}
