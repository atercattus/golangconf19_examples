package scenes

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	engoCommon "github.com/EngoEngine/engo/common"
	"github.com/atercattus/golangconf19_examples/engo/common"
	"github.com/atercattus/golangconf19_examples/engo/systems"
	"image/color"
	"log"
	"math"
)

type (
	GameScene struct {
		common.CommonScene

		road *systems.RoadSystem

		player *systems.GopherSystem
	}
)

const (
	imgFilenameRunBtn  = `images/run.png`
	imgFilenameRoad    = `images/tunnel_road.jpg`
	imgFilenameScooter = `images/scooter.png`
	imgFilenameGopher  = `images/gopher.png`

	// начало дороги по вертикали от верхней части экрана
	roadMinY = 50
	// высота обочины
	topH = 23
	// высота полосы дороги
	roadH = 115
)

func (scene *GameScene) Preload() {
	scene.CommonScene.Preload()

	_ = engo.Files.Load(imgFilenameRunBtn)
	_ = engo.Files.Load(imgFilenameRoad)
	_ = engo.Files.Load(imgFilenameScooter)
	_ = engo.Files.Load(imgFilenameGopher)
}

func (scene *GameScene) Setup(u engo.Updater) {
	scene.CommonScene.Setup(u)
	engoCommon.SetBackground(color.RGBA{80, 80, 80, 255})
	world := scene.World()

	//fmt.Println(`USERNAME:`, common.GlobalData.UserName)

	btn := scene.SetupRunBtn(u)
	scene.MakeRoadLine(u)

	mouseSystem := &engoCommon.MouseSystem{}
	world.AddSystem(mouseSystem)
	mouseSystem.Add(&btn.BasicEntity, &btn.MouseComponent, &btn.SpaceComponent, &btn.RenderComponent)

	scene.player = systems.NewGopherSystem(scene.World())
	scene.player.MoveTo(engo.Point{
		X: 0,
		Y: roadMinY + topH - scene.player.Height()/2,
	})
	//gopher.SetRelativeSpeed(0)
	//gopher.SetIsDriving(true)
}

func (scene *GameScene) SetupRunBtn(u engo.Updater) *common.Sprite {
	texture, err := engoCommon.LoadedSprite(imgFilenameRunBtn)
	if err != nil {
		log.Println("Unable to load texture: " + err.Error())
		return nil
	}

	runBtn := common.NewSprite(texture, &common.SpriteOptions{
		Position: engo.Point{
			X: engo.WindowWidth() - texture.Width(),
			Y: engo.WindowHeight() - texture.Height(),
		},
		Width:  texture.Width(),
		Height: texture.Height(),
		ZIndex: 1,
		Shader: engoCommon.HUDShader,
	})
	runBtn.AddToWorld(u)

	scene.World().AddSystem(&systems.MousableSystem{Callback: func(dt float32) {
		scale := runBtn.Scale.X
		if runBtn.MouseComponent.Clicked {
			scale += 10 * dt
			scene.road.Speed += 8
			scene.player.SetRelativeSpeed(scene.player.GetRelativeSpeed() + 2)
			scene.player.SetIsDriving(true)
		} else {
			scale -= 1 * dt
		}

		if scale >= 1 && scale < 1.5 {
			runBtn.Scale.X = scale
			runBtn.Scale.Y = scale
		}
	}})

	text := common.NewText(`RUN`, common.Font80, &common.TextOptions{
		SpriteOptions: common.SpriteOptions{
			Scale:    0.5,
			Position: runBtn.SpaceComponent.Position,
			ZIndex:   2,
			Width:    runBtn.Width,
			Height:   runBtn.Height,
		},
		AlignH: common.TextAlignCenter,
		AlignV: common.TextAlignCenter,
	})
	text.AddToWorld(u)

	return runBtn
}

func (scene *GameScene) MakeRoadLine(u engo.Updater) {
	const spriteH = 512
	const W = 512
	sheet := engoCommon.NewAsymmetricSpritesheetFromFile(imgFilenameRoad, []engoCommon.SpriteRegion{
		{Position: engo.Point{X: 0, Y: 0}, Width: W, Height: topH},              // top
		{Position: engo.Point{X: 0, Y: topH}, Width: W, Height: roadH},          // line
		{Position: engo.Point{X: 0, Y: spriteH - topH}, Width: W, Height: topH}, // bottom
	})

	// сколько отрезков идет в стык по горизонтали
	columns := 1 + int(math.Ceil(float64(engo.WindowWidth())/W))
	// число игроков
	const roads = 1

	scene.road = &systems.RoadSystem{
		Speed: 0,
		MinX:  -W,
		MaxX:  W * float32(columns),
	}
	world, _ := u.(*ecs.World)
	world.AddSystem(scene.road)

	for column := 0; column < columns; column++ {
		x := float32(column * W)

		if txtTop := sheet.Drawable(0).(engoCommon.Texture); true {
			roadLine := common.NewSprite(&txtTop, &common.SpriteOptions{
				Position: engo.Point{X: x, Y: roadMinY},
			})
			scene.road.Add(roadLine)
			roadLine.AddToWorld(u)
		}

		txtRoad := sheet.Drawable(1).(engoCommon.Texture)
		for i := 0; i < roads; i++ {
			roadLine := common.NewSprite(&txtRoad, &common.SpriteOptions{
				Position: engo.Point{X: x, Y: float32(roadMinY + topH + i*roadH)},
			})
			scene.road.Add(roadLine)
			roadLine.AddToWorld(u)
		}

		if txtBottom := sheet.Drawable(2).(engoCommon.Texture); true {
			roadLine := common.NewSprite(&txtBottom, &common.SpriteOptions{
				Position: engo.Point{X: x, Y: float32(roadMinY + topH + roads*roadH - topH)}, // -topH чтобы спрятать последнюю разметку
			})
			scene.road.Add(roadLine)
			roadLine.AddToWorld(u)
		}
	}
}
