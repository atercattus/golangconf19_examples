package render

import (
	"strconv"
	"syscall/js"
)

type (
	Text2D struct {
		tr *TextRenderer

		text         string
		pos          Point
		size         int
		font         string
		alignX       string // https://www.w3schools.com/tags/canvas_textalign.asp
		outlineWidth int
	}

	TextRenderer struct {
		jsCanvas js.Value
		textCtx  js.Value

		canvasWidth  float32
		canvasHeight float32

		texts map[*Text2D]bool

		changed bool
	}
)

func NewText2D(tr *TextRenderer, text string, size int, pos Point) *Text2D {
	t := &Text2D{
		tr:     tr,
		text:   text,
		pos:    pos,
		size:   size,
		font:   strconv.Itoa(size) + `px Roboto`,
		alignX: `center`,
	}

	return t
}

func (t *Text2D) Pos() Point {
	return t.pos
}

func (t *Text2D) Size() float32 {
	return float32(t.size)
}

func (t *Text2D) AlignX() string {
	return t.alignX
}

func (t *Text2D) SetAlignX(align string) {
	t.alignX = align
	t.tr.markChanged()
}

func (t *Text2D) SetOutlineWidth(width int) {
	t.outlineWidth = width
	t.tr.markChanged()
}

func (t *Text2D) SetText(text string) {
	t.text = text
	t.tr.markChanged()
}

func NewTextRenderer(canvasId string) *TextRenderer {
	tr := &TextRenderer{
		texts: make(map[*Text2D]bool),
	}

	tr.jsCanvas = js.Global().Get(`document`).Call(`getElementById`, canvasId)

	tr.textCtx = tr.jsCanvas.Call(`getContext`, `2d`)

	tr.Resize()

	return tr
}

func (tr *TextRenderer) markChanged() {
	tr.changed = true
}

func (tr *TextRenderer) Resize() {
	tr.canvasWidth = float32(tr.jsCanvas.Get(`width`).Float())
	tr.canvasHeight = float32(tr.jsCanvas.Get(`height`).Float())

	//tr.textCtx.Set(`textAlign`, `center`)
	tr.textCtx.Set(`textBaseline`, `middle`)

	tr.markChanged()
}

func (tr *TextRenderer) AddText(text string, size int, pos Point) *Text2D {
	txt := NewText2D(tr, text, size, pos)
	tr.texts[txt] = true
	tr.markChanged()
	return txt
}

func (tr *TextRenderer) DeleteText(text *Text2D) {
	delete(tr.texts, text)
	tr.markChanged()
}

func (tr *TextRenderer) Draw() {
	// нужно перерисовывать только если поменялось
	if !tr.changed {
		return
	}
	tr.changed = false

	//println(`TEXT REPAINT`)

	tr.textCtx.Call(`clearRect`, 0, 0, tr.canvasWidth, tr.canvasHeight)

	for text, _ := range tr.texts {
		tr.textCtx.Set(`font`, text.font)
		tr.textCtx.Set(`textAlign`, text.alignX)

		if text.outlineWidth > 0 {
			tr.textCtx.Call(`save`)
			tr.textCtx.Set(`strokeStyle`, `black`)
			tr.textCtx.Set(`fillStyle`, `white`)
			tr.textCtx.Set(`lineWidth`, text.outlineWidth)
			tr.textCtx.Call(`strokeText`, text.text, text.pos.X, text.pos.Y)
			tr.textCtx.Call(`fillText`, text.text, text.pos.X, text.pos.Y)
			tr.textCtx.Call(`restore`)
		} else {
			tr.textCtx.Set(`fillStyle`, `black`)
			tr.textCtx.Call(`fillText`, text.text, text.pos.X, text.pos.Y)
		}
	}
}
