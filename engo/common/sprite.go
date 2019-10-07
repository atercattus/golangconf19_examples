package common

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type (
	Sprite struct {
		ecs.BasicEntity
		common.RenderComponent
		common.SpaceComponent
		common.MouseComponent
	}

	SpriteOptions struct {
		Position engo.Point
		ZIndex   float32
		Width    float32
		Height   float32
		Rotation float32
		Scale    float32

		Repeat    common.TextureRepeating
		MinFilter common.ZoomFilter
		MagFilter common.ZoomFilter
		Shader    common.Shader
	}
)

func NewSprite(texture *common.Texture, options *SpriteOptions) *Sprite {
	sprite := &Sprite{}

	sprite.BasicEntity = ecs.NewBasic()
	sprite.SpaceComponent = common.SpaceComponent{
		Width:  texture.Width(),
		Height: texture.Height(),
	}

	sprite.RenderComponent = common.RenderComponent{
		Drawable: texture,
	}

	if options != nil {
		sprite.applyOptions(options)
	}

	return sprite
}

func (s *Sprite) applyOptions(options *SpriteOptions) {
	if options.Scale != 0 {
		scale := options.Scale
		s.RenderComponent.Scale = engo.Point{X: scale, Y: scale}
	}
	s.RenderComponent.SetZIndex(options.ZIndex)
	s.RenderComponent.SetMinFilter(options.MinFilter)
	if options.Shader != nil {
		s.RenderComponent.SetShader(options.Shader)
	}
	s.RenderComponent.Repeat = options.Repeat

	s.SpaceComponent.Position = options.Position
	s.SpaceComponent.Width = options.Width
	s.SpaceComponent.Height = options.Height
	s.SpaceComponent.Rotation = options.Rotation
}

func (s *Sprite) AddToWorld(u engo.Updater) {
	world, _ := u.(*ecs.World)

	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&s.BasicEntity, &s.RenderComponent, &s.SpaceComponent)
			return // ?
		}
	}
}
