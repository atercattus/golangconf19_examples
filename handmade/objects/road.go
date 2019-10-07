package objects

import (
	"github.com/atercattus/golangconf19_examples/handmade/events"
	"github.com/atercattus/golangconf19_examples/handmade/render"
	"github.com/atercattus/golangconf19_examples/handmade/resources"
	"math"
)

type (
	Road struct {
		batchRoad        *render.DrawCallBatch
		columnsRoadCount int

		batchBg        *render.DrawCallBatch
		columnsBgCount int

		topY                       float32
		bottomY                    float32
		middleLineHeightRenderInPx float32

		roadSprites []*render.Sprite
		bgSprites   []*render.Sprite

		playersNames []*render.Text2D

		speed float32

		batchFinishLine *render.DrawCallBatch
		finishLine      *render.Sprite
	}
)

const (
	// высота верхнего и нижнего бордюров в пикселях
	roadTopLineHeightInPx    = 22.0
	roadBottomLineHeightInPx = 22.0
	// высота основной проезжей части в пикселях
	roadMiddleLineHeightInPx = 239.0
)

func NewRoad(renderer *render.WebGLRender, players []events.PlayerInfo, maxPlayers int, topY float32) *Road {
	road := &Road{}

	road.initLines(renderer, players, maxPlayers, topY)
	road.initBackground(renderer)

	playersCount := len(players)
	if playersCount == 0 {
		// если выводим дорогу при сборе участников, то там нет гоферов (?)
		playersCount = maxPlayers
	}
	road.initFinishLine(renderer, playersCount)

	return road
}

func (road *Road) initFinishLine(renderer *render.WebGLRender, playersCount int) {
	road.batchFinishLine = render.NewDrawCallBatch(renderer, resources.TextureFinishLine)
	road.finishLine = road.batchFinishLine.AddSprite()
	road.finishLine.SetColor(1, 1, 1, 0.6)
	t := road.topY
	b := road.bottomY
	road.finishLine.SizeTo(render.Point{X: 32, Y: b - t})
	road.finishLine.MoveTo(render.Point{X: -road.finishLine.Size().X, Y: t + road.finishLine.Size().Y/2})
}

func (road *Road) initBackground(renderer *render.WebGLRender) {
	texture := resources.TextureBkg
	size := texture.Size()

	road.batchBg = render.NewDrawCallBatch(renderer, texture)

	W, _ := renderer.Size()

	road.columnsBgCount = 1 + int(math.Ceil(float64(W/texture.Size().X)))

	for column := 0; column < road.columnsBgCount; column++ {
		x := size.X/2 + float32(column)*size.X

		sprite := road.batchBg.AddSprite()
		road.bgSprites = append(road.bgSprites, sprite)
		sprite.MoveTo(render.Point{X: x, Y: size.Y * 0.3})
	}
}

func (road *Road) initLines(renderer *render.WebGLRender, players []events.PlayerInfo, maxPlayers int, topY float32,
) {
	texture := resources.TextureRoad

	road.batchRoad = render.NewDrawCallBatch(renderer, texture)

	W, H := renderer.Size()
	road.topY = topY

	road.speed = 0

	// высота основной проезжей части в пикселях при отрисовке (подстраиваясь под высоту экрана)
	road.middleLineHeightRenderInPx = (H - road.topY) * 0.75 / float32(maxPlayers)

	playersCount := len(players)
	if playersCount == 0 {
		// если выводим дорогу при сборе участников, то там нет гоферов (?)
		playersCount = maxPlayers
	}

	// позиция нижнего бордюра
	road.bottomY = road.topY + roadTopLineHeightInPx + float32(playersCount)*road.middleLineHeightRenderInPx

	// сколько отрезков идет в стык по горизонтали
	road.columnsRoadCount = 1 + int(math.Ceil(float64(W/texture.Size().X)))

	// верхний бордюр
	road.addTop(texture)

	// полосы дороги
	for roadIdx := 0; roadIdx < playersCount; roadIdx++ {
		const fontSize = 40
		y := road.addLine(texture, roadIdx) - road.middleLineHeightRenderInPx/2 + fontSize/2

		txt := renderer.AddText(``, fontSize, render.Point{X: W * 0.99, Y: y})
		txt.SetAlignX(`end`)
		txt.SetOutlineWidth(8)
		road.playersNames = append(road.playersNames, txt)

		if len(players) > roadIdx {
			txt.SetText(players[roadIdx].Name)
		}
	}

	// нижний бордюр
	road.addBottom(texture)
}

func (road *Road) addTop(fullSprite *render.Texture) {
	const h = roadTopLineHeightInPx / 512.0

	size := fullSprite.Size()

	for column := 0; column < road.columnsRoadCount; column++ {
		x := size.X/2 + float32(column)*size.X
		y := road.topY + roadTopLineHeightInPx/2

		sprite := road.batchRoad.AddSprite()
		road.roadSprites = append(road.roadSprites, sprite)
		sprite.MoveTo(render.Point{X: x, Y: y})

		// растягиваю спрайт и текстурные координаты так, чтобы визуально отрезать от всей картинки верхнюю полоску
		sprite.SizeTo(render.Point{X: size.X, Y: roadTopLineHeightInPx})
		sprite.SetTexCoords([render.VerticesTextureSpriteSize]float32{
			0, h,
			0, 0,
			1, 0,
			1, h,
		})
	}
}

func (road *Road) addBottom(fullSprite *render.Texture) {
	const h = roadBottomLineHeightInPx / 512.0

	size := fullSprite.Size()

	for column := 0; column < road.columnsRoadCount; column++ {
		x := size.X/2 + float32(column)*size.X
		y := road.bottomY - roadBottomLineHeightInPx/2 // закрываю нижним бордюром последнюю разделительную разметку

		sprite := road.batchRoad.AddSprite()
		road.roadSprites = append(road.roadSprites, sprite)
		sprite.MoveTo(render.Point{X: x, Y: y})

		// растягиваю спрайт и текстурные координаты так, чтобы визуально отрезать от всей картинки нижнюю полоску
		sprite.SizeTo(render.Point{X: size.X, Y: roadBottomLineHeightInPx})
		sprite.SetTexCoords([render.VerticesTextureSpriteSize]float32{
			0, 1,
			0, 1 - h,
			1, 1 - h,
			1, 1,
		})
	}
}

func (road *Road) addLine(fullSprite *render.Texture, lineIdx int) float32 {
	const topH = roadTopLineHeightInPx / 512.0
	const midH = roadMiddleLineHeightInPx / 512.0

	size := fullSprite.Size()

	y := road.GetLineY(lineIdx)

	for column := 0; column < road.columnsRoadCount; column++ {
		x := size.X/2 + float32(column)*size.X

		sprite := road.batchRoad.AddSprite()
		road.roadSprites = append(road.roadSprites, sprite)
		sprite.MoveTo(render.Point{X: x, Y: y})

		// растягиваю спрайт и текстурные координаты так, чтобы визуально отрезать от всей картинки середину
		sprite.SizeTo(render.Point{X: size.X, Y: road.middleLineHeightRenderInPx})
		sprite.SetTexCoords([render.VerticesTextureSpriteSize]float32{
			0, topH + midH,
			0, topH,
			1, topH,
			1, topH + midH,
		})
	}

	return y
}

func (road *Road) move(dt float32) {
	if road.speed == 0 {
		return
	}

	delta := render.Point{X: -road.speed * dt}
	road.moveLines(delta)
	road.moveBgs(delta)
	road.finishLine.MoveBy(delta)
}

func (road *Road) moveLines(delta render.Point) {
	for _, sprite := range road.roadSprites {
		sprite.MoveBy(delta)
		if sprite.Pos().X < -sprite.Size().X/2 {
			sprite.MoveBy(render.Point{X: sprite.Size().X * float32(road.columnsRoadCount)})
		}
	}
}

func (road *Road) moveBgs(delta render.Point) {
	for _, sprite := range road.bgSprites {
		sprite.MoveBy(delta)
		if sprite.Pos().X < -sprite.Size().X/2 {
			sprite.MoveBy(render.Point{X: sprite.Size().X * float32(road.columnsBgCount)})
		}
	}
}

func (road *Road) GetLineY(lineIdx int) float32 {
	return road.topY +
		roadTopLineHeightInPx +
		road.middleLineHeightRenderInPx/2 +
		float32(lineIdx)*road.middleLineHeightRenderInPx
}

func (road *Road) GetLineHeightInPx() float32 {
	return road.middleLineHeightRenderInPx
}

func (road *Road) SetPlayerName(idx int, name string) {
	if idx < len(road.playersNames) {
		road.playersNames[idx].SetText(name)
	}
}

func (road *Road) Delete(renderer *render.WebGLRender) {
	for _, txt := range road.playersNames {
		renderer.DeleteText(txt)
	}
}

func (road *Road) Draw(dt float32) {
	road.batchBg.Draw()
	road.batchRoad.Draw()
	road.batchFinishLine.Draw()
	road.move(dt)
}

func (road *Road) SetSpeed(speed float32) {
	road.speed = speed
}

func (road *Road) SetFinishDistance(distance float32) {
	pos := road.finishLine.Pos()
	pos.X = distance
	road.finishLine.MoveTo(pos)
}
