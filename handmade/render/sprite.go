package render

import "math"

type (
	Sprite struct {
		batch       *DrawCallBatch
		texture     *Texture
		texCoords   [VerticesTextureSpriteSize]float32
		colors      [VerticesColorsSpriteSize]float32
		verticesIdx int

		pos     Point
		size    Point
		angle   float32 // градусы
		visible bool
	}
)

func NewSprite(batch *DrawCallBatch, texture *Texture, verticesIdx int) *Sprite {
	sprite := &Sprite{
		batch:       batch,
		texture:     texture,
		verticesIdx: verticesIdx,

		size:    texture.size,
		visible: true,
	}

	copy(sprite.texCoords[:], texCoordsEtalon[:])
	copy(sprite.colors[:], colorsEtalon[:])

	sprite.updateVertices()

	return sprite
}

func (s *Sprite) updateVertices() {
	s.batch.updateVertices(s.verticesIdx, s.pos, s.size, s.angle)
}

func (s *Sprite) updateTexCoords() {
	s.batch.updateTexCoords(s.verticesIdx, s.texCoords)
}

func (s *Sprite) updateColors() {
	s.batch.updateColors(s.verticesIdx, s.colors)
}

func (s *Sprite) Pos() Point {
	return s.pos
}

func (s *Sprite) Size() Point {
	return s.size
}

func (s *Sprite) Angle() float32 {
	return s.angle
}

func (s *Sprite) MoveTo(pos Point) {
	s.pos = pos
	s.updateVertices()
}

func (s *Sprite) MoveBy(delta Point) {
	s.pos.X += delta.X
	s.pos.Y += delta.Y
	s.updateVertices()
}

func (s *Sprite) SizeTo(size Point) {
	if size.Y == 0 && size.X == 0 {
		return
	}

	// Если один из размеров не задан, но просто масштабирую пропорционально
	if size.Y == 0 {
		size.Y = size.X / s.size.X * s.size.Y
	} else if size.X == 0 {
		size.X = size.Y / s.size.Y * s.size.X
	}

	s.size = size
	s.updateVertices()
}

func (s *Sprite) RotateTo(angle float32) {
	for angle < -360 {
		angle += 360
	}
	for angle > 360 {
		angle -= 360
	}

	s.angle = angle
	s.updateVertices()
}

func (s *Sprite) SetTexCoords(texCoords [VerticesTextureSpriteSize]float32) {
	s.texCoords = texCoords
	s.updateTexCoords()
}

func (s *Sprite) SetColors(colors [VerticesColorsSpriteSize]float32) {
	s.colors = colors
	s.updateColors()
}

func (s *Sprite) SetColor(r, g, b, a float32) {
	s.SetColors([VerticesColorsSpriteSize]float32{
		r, g, b, a,
		r, g, b, a,
		r, g, b, a,
		r, g, b, a,
	})
}

func (s *Sprite) IsPointInside(pos Point) bool {
	// оч тупо
	return math.Abs(float64(pos.X-s.pos.X)) <= float64(s.size.X) &&
		math.Abs(float64(pos.Y-s.pos.Y)) <= float64(s.size.Y)
}

//func (s *Sprite) Show(show bool) {
//	s.visible = show
//	// ToDo:
//}
//
//func (s *Sprite) Delete() {
//	// ToDo:
//}
