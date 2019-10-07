package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	engoCommon "github.com/EngoEngine/engo/common"
	"github.com/atercattus/golangconf19_examples/engo/common"
	"math"
)

type (
	GopherSystemUpdate func(dt float32, gopher *GopherSystem)

	GopherSystem struct {
		world         *ecs.World
		gopherSprite  *common.Sprite
		scooterSprite *common.Sprite

		updateCb      GopherSystemUpdate
		relativeSpeed float32
		driveTime     float32
		isDriving     bool
	}
)

var (
	_ ecs.System = &GopherSystem{}

	textureGopher  *engoCommon.Texture
	textureScooter *engoCommon.Texture
)

const (
	imgFilenameScooter = `images/scooter.png`
	imgFilenameGopher  = `images/gopher.png`
)

func NewGopherSystem(world *ecs.World) (sys *GopherSystem) {
	sys = &GopherSystem{
		world: world,
	}
	sys.loadSprites()

	world.AddSystem(sys)

	return sys
}

func (gopher *GopherSystem) SetUpdateCb(updateCb GopherSystemUpdate) {
	gopher.updateCb = updateCb
}

func (gopher *GopherSystem) SetRelativeSpeed(speed float32) {
	gopher.relativeSpeed = speed
}

func (gopher *GopherSystem) GetRelativeSpeed() float32 {
	return gopher.relativeSpeed
}

func (gopher *GopherSystem) SetIsDriving(isDriving bool) {
	gopher.isDriving = isDriving
}

func (gopher *GopherSystem) IsDriving() bool {
	return gopher.isDriving
}

func (*GopherSystem) Remove(ecs.BasicEntity) {
}

func (gopher *GopherSystem) Update(dt float32) {
	if gopher.isDriving || (gopher.relativeSpeed != 0) {
		pos := gopher.Position()

		if gopher.isDriving {
			gopher.driveTime += dt

			// небольшой сдвиг по вертикали для реализма
			pos.Y += float32(math.Sin(float64(gopher.driveTime)) / 6.0)
		}

		pos.X += gopher.relativeSpeed * dt

		gopher.MoveTo(pos)
	}

	// пользовательский код
	if gopher.updateCb != nil {
		gopher.updateCb(dt, gopher)
	}
}

func (gopher *GopherSystem) Position() engo.Point {
	return gopher.gopherSprite.Position
}

func (gopher *GopherSystem) Width() float32 {
	return gopher.gopherSprite.Width * gopher.gopherSprite.Scale.X // ToDo: заменить на объединение обоих AABB
}

func (gopher *GopherSystem) Height() float32 {
	return gopher.gopherSprite.Height * gopher.gopherSprite.Scale.Y // ToDo: заменить на объединение обоих AABB
}

func (gopher *GopherSystem) MoveTo(pos engo.Point) {
	gopher.gopherSprite.Position = pos

	gopher.scooterSprite.Position.X = pos.X + 20
	gopher.scooterSprite.Position.Y = pos.Y + 15
}

func (gopher *GopherSystem) loadSprites() {
	var err error

	_ = engo.Files.Load(imgFilenameGopher)
	_ = engo.Files.Load(imgFilenameScooter)

	if textureGopher == nil {
		if textureGopher, err = engoCommon.LoadedSprite(imgFilenameGopher); err != nil {
			panic(err)
		}
	}
	if textureScooter == nil {
		if textureScooter, err = engoCommon.LoadedSprite(imgFilenameScooter); err != nil {
			panic(err)
		}
	}

	gopher.gopherSprite = common.NewSprite(textureGopher, &common.SpriteOptions{
		Width:  textureGopher.Width(),
		Height: textureGopher.Height(),
		Scale:  0.4,
		ZIndex: 2,
	})
	gopher.gopherSprite.AddToWorld(gopher.world)

	gopher.scooterSprite = common.NewSprite(textureScooter, &common.SpriteOptions{
		Width:  textureScooter.Width(),
		Height: textureScooter.Height(),
		Scale:  0.35,
		ZIndex: 3,
	})
	gopher.scooterSprite.AddToWorld(gopher.world)
}
