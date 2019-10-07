package common

import (
	"bytes"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	engoCommon "github.com/EngoEngine/engo/common"
	"golang.org/x/image/font/gofont/gosmallcaps"
	"image/color"
)

type (
	TextAlign byte

	TextOptions struct {
		SpriteOptions
		AlignH TextAlign
		AlignV TextAlign
	}

	Text struct {
		Sprite

		AlignH TextAlign
		AlignV TextAlign
	}
)

const (
	TextAlignMin = TextAlign(iota)
	TextAlignCenter
	TextAlignMax
)

var (
	Font80     *engoCommon.Font
	Font80Size = 80

	Font80DarkBlue *engoCommon.Font
)

const (
	fontFileName = `GoSmallcaps.ttf`
)

func FontsPreload() {
	_ = engo.Files.LoadReaderData(fontFileName, bytes.NewReader(gosmallcaps.TTF))
	Font80 = &engoCommon.Font{
		URL:  fontFileName,
		FG:   color.Black,
		Size: float64(Font80Size),
	}
	_ = Font80.CreatePreloaded()

	_ = engo.Files.LoadReaderData(fontFileName, bytes.NewReader(gosmallcaps.TTF))
	Font80DarkBlue = &engoCommon.Font{
		URL:  fontFileName,
		FG:   color.RGBA{R: 60, G: 60, B: 150, A: 255},
		Size: float64(Font80Size),
	}
	_ = Font80DarkBlue.CreatePreloaded()
}

func NewText(text string, font *engoCommon.Font, options *TextOptions) *Text {
	textSprite := Text{}
	textSprite.BasicEntity = ecs.NewBasic()
	txt := engoCommon.Text{
		Font: font,
		Text: text,
	}
	textSprite.RenderComponent = engoCommon.RenderComponent{
		Drawable: txt,
	}

	textSprite.SetShader(engoCommon.TextHUDShader)
	textSprite.SpaceComponent = engoCommon.SpaceComponent{}

	if options != nil {
		textSprite.applyOptions(options)
	}

	return &textSprite
}

func (t *Text) applyOptions(options *TextOptions) {
	t.Sprite.applyOptions(&options.SpriteOptions)

	t.AlignH = options.AlignH
	t.AlignV = options.AlignV
	t.Realign()
}

func (t *Text) Realign() {
	txt := t.RenderComponent.Drawable.(engoCommon.Text)
	scale := t.RenderComponent.Scale.X
	if scale == 0 {
		scale = 1
	}

	wDiff := t.SpaceComponent.Width - (txt.Width() * scale)
	switch t.AlignH {
	case TextAlignCenter:
		t.SpaceComponent.Position.X += wDiff / 2
	case TextAlignMax:
		t.SpaceComponent.Position.X += wDiff
	}

	hDiff := t.SpaceComponent.Height - (txt.Height() * scale)
	switch t.AlignV {
	case TextAlignCenter:
		t.SpaceComponent.Position.Y += hDiff / 2
	case TextAlignMax:
		t.SpaceComponent.Position.Y += hDiff
	}
}
