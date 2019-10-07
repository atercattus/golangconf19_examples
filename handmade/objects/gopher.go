package objects

import (
	"github.com/atercattus/golangconf19_examples/handmade/render"
	"github.com/atercattus/golangconf19_examples/handmade/resources"
)

type (
	Gopher struct {
		batchGopher *render.DrawCallBatch
		batchPaw    *render.DrawCallBatch

		gopher *render.Sprite
		paw    *render.Sprite
		pawPos render.Point
	}
)

func NewGopher(renderer *render.WebGLRender) *Gopher {
	gopher := &Gopher{}

	gopher.batchGopher = render.NewDrawCallBatch(renderer, resources.TextureGopher)
	gopher.batchPaw = render.NewDrawCallBatch(renderer, resources.TextureGopherPaw)

	gopher.gopher = gopher.batchGopher.AddSprite()
	gopher.paw = gopher.batchPaw.AddSprite()

	return gopher
}

func (gopher *Gopher) Size() render.Point {
	return gopher.gopher.Size()
}

func (gopher *Gopher) MoveTo(pos render.Point) {
	cur := gopher.gopher.Pos()
	delta := render.Point{
		X: pos.X - cur.X,
		Y: pos.Y - cur.Y,
	}
	gopher.MoveBy(delta)
}

func (gopher *Gopher) MoveBy(delta render.Point) {
	gopher.gopher.MoveBy(delta)
	gopher.paw.MoveBy(delta)
	gopher.pawPos = gopher.paw.Pos()
}

func (gopher *Gopher) SizeTo(size render.Point) {
	gopher.gopher.SizeTo(size)
	gopher.paw.SizeTo(size)
}

func (gopher *Gopher) Pos() render.Point {
	return gopher.gopher.Pos()
}

func (gopher *Gopher) Draw() {
	gopher.batchGopher.Draw()
	gopher.batchPaw.Draw()
}
