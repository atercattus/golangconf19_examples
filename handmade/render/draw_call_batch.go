package render

import (
	"github.com/nuberu/webgl"
	"github.com/nuberu/webgl/types"
	"math"
)

type (
	DrawCallBatch struct {
		renderer *WebGLRender
		texture  *Texture

		indices   []uint16
		vertices  []float32
		texCoords []float32
		colors    []float32

		verticesBuffer  *types.Buffer
		texCoordsBuffer *types.Buffer
		colorsBuffer    *types.Buffer
		indexBuffer     *types.Buffer

		verticesAreModified  bool
		texCoordsAreModified bool
		colorsAreModified    bool

		shader *types.Program

		maxSpritesCount int
		lastUsedIdx     int
		spritesCount    int
	}
)

var (
	indicesEtalon = [SpriteIndicesCount]uint16{
		0, 1, 3, 1, 2, 3,
	}

	texCoordsEtalon = [VerticesTextureSpriteSize]float32{
		0, 1, 0, 0, 1, 0, 1, 1,
	}

	colorsEtalon = [VerticesColorsSpriteSize]float32{
		1, 1, 1, 1, // r
		1, 1, 1, 1, // g
		1, 1, 1, 1, // b
		1, 1, 1, 1, // a
	}
)

func NewDrawCallBatch(renderer *WebGLRender, texture *Texture) *DrawCallBatch {
	batch := &DrawCallBatch{
		renderer:        renderer,
		texture:         texture,
		maxSpritesCount: 200, // ToDo:
		lastUsedIdx:     0,   // резервирую 0 элемент под пустышку
		spritesCount:    1,   // резервирую 0 элемент под пустышку
	}

	batch.precalcVertices()
	batch.precalcIndices()
	batch.precalcTexCoord()
	batch.precalcColors()
	batch.setupShaders()

	return batch
}

func (batch *DrawCallBatch) precalcVertices() {
	cnt := batch.maxSpritesCount
	batch.vertices = make([]float32, cnt*VerticesSpriteSize)

	batch.verticesBuffer = batch.renderer.CreateArrayBuffer(batch.vertices)
}

func (batch *DrawCallBatch) precalcIndices() {
	cnt := batch.maxSpritesCount
	batch.indices = make([]uint16, cnt*SpriteIndicesCount)
	for sprite := 0; sprite < cnt; sprite++ {
		dstIdx := sprite * SpriteIndicesCount
		offset := sprite * SpriteRectIndices
		for i, base := range indicesEtalon {
			batch.indices[dstIdx+i] = base + uint16(offset)
		}
	}

	batch.indexBuffer = batch.renderer.CreateElementArrayBufferUI16(batch.indices)
}

func (batch *DrawCallBatch) precalcTexCoord() {
	cnt := batch.maxSpritesCount
	batch.texCoords = make([]float32, cnt*VerticesTextureSpriteSize)
	for sprite := 0; sprite < cnt; sprite++ {
		copy(batch.texCoords[sprite*VerticesTextureSpriteSize:], texCoordsEtalon[:])
	}

	batch.texCoordsBuffer = batch.renderer.CreateArrayBuffer(batch.texCoords)
}

func (batch *DrawCallBatch) precalcColors() {
	cnt := batch.maxSpritesCount
	batch.colors = make([]float32, cnt*VerticesColorsSpriteSize)
	for sprite := 0; sprite < cnt; sprite++ {
		copy(batch.colors[sprite*VerticesColorsSpriteSize:], colorsEtalon[:])
	}

	batch.colorsBuffer = batch.renderer.CreateArrayBuffer(batch.colors)
}

func (batch *DrawCallBatch) setupShaders() {
	// ToDo: кешить одинаковые шейдеры
	batch.shader = batch.renderer.LoadShaderProgram(vertexShaderCode, fragmentShaderCode)
}

func (batch *DrawCallBatch) bindShaders() {
	gl := batch.renderer.GetGlCtx()
	renderer := batch.renderer

	gl.UseProgram(batch.shader)

	renderer.SetupPosVertexBufferAttrib(batch.shader, batch.verticesBuffer)
	renderer.SetupTextureVertexBufferAttrib(batch.shader, batch.texCoordsBuffer)
	renderer.SetupColorsVertexBufferAttrib(batch.shader, batch.colorsBuffer)
	renderer.SetupPVUniform(batch.shader)
	renderer.SetupSamplerUniform(batch.shader, 0) // webgl.TEXTURE0 + texIdx
}

func (batch *DrawCallBatch) Bind() {
	gl := batch.renderer.GetGlCtx()

	gl.BindBuffer(webgl.ARRAY_BUFFER, batch.verticesBuffer)
	gl.BindBuffer(webgl.ARRAY_BUFFER, batch.texCoordsBuffer)
	gl.BindBuffer(webgl.ARRAY_BUFFER, batch.colorsBuffer)
	gl.BindBuffer(webgl.ELEMENT_ARRAY_BUFFER, batch.indexBuffer)

	batch.bindShaders()
}

func (batch *DrawCallBatch) updateVertices(spriteIdx int, pos Point, size Point, angle float32) {
	idx := spriteIdx * VerticesSpriteSize

	vertices := batch.vertices

	if lastIdx := idx + VerticesSpriteSize; lastIdx > len(vertices) {
		println(`wrong idx`)
		return
	}

	l := pos.X - size.X/2
	r := pos.X + size.X/2
	t := pos.Y - size.Y/2
	b := pos.Y + size.Y/2

	x1, y1 := l, b
	x2, y2 := l, t
	x3, y3 := r, t
	x4, y4 := r, b

	sin64, cos64 := math.Sincos(float64(angle) / 180 * math.Pi)
	sin, cos := float32(sin64), float32(cos64)

	rot := func(px, py float32) (x float32, y float32) {
		x = cos*(px-pos.X) - sin*(py-pos.Y) + pos.X
		y = sin*(px-pos.X) + cos*(py-pos.Y) + pos.Y
		return
	}

	vertices[idx+0], vertices[idx+1] = rot(x1, y1)
	vertices[idx+2], vertices[idx+3] = rot(x2, y2)
	vertices[idx+4], vertices[idx+5] = rot(x3, y3)
	vertices[idx+6], vertices[idx+7] = rot(x4, y4)

	batch.verticesAreModified = true
}

func (batch *DrawCallBatch) updateTexCoords(spriteIdx int, texCoords [VerticesTextureSpriteSize]float32) {
	idx := spriteIdx * VerticesTextureSpriteSize

	if lastIdx := idx + VerticesTextureSpriteSize; lastIdx > len(batch.texCoords) {
		println(`wrong idx`)
		return
	}

	copy(batch.texCoords[idx:], texCoords[:])

	batch.texCoordsAreModified = true
}

func (batch *DrawCallBatch) updateColors(spriteIdx int, colors [VerticesColorsSpriteSize]float32) {
	idx := spriteIdx * VerticesColorsSpriteSize

	if lastIdx := idx + VerticesColorsSpriteSize; lastIdx > len(batch.colors) {
		println(`wrong idx`)
		return
	}

	copy(batch.colors[idx:], colors[:])

	batch.colorsAreModified = true
}

func (batch *DrawCallBatch) Draw() {
	if batch.spritesCount == 0 {
		return
	}

	gl := batch.renderer.GetGlCtx()

	if batch.verticesAreModified {
		batch.renderer.UpdateArrayBuffer(batch.verticesBuffer, batch.vertices)
		batch.verticesAreModified = false
	}

	if batch.texCoordsAreModified {
		batch.renderer.UpdateArrayBuffer(batch.texCoordsBuffer, batch.texCoords)
		batch.texCoordsAreModified = false
	}

	if batch.colorsAreModified {
		batch.renderer.UpdateArrayBuffer(batch.colorsBuffer, batch.colors)
		batch.colorsAreModified = false
	}

	batch.Bind()
	batch.texture.Bind()
	gl.DrawElements(webgl.TRIANGLES, batch.spritesCount*SpriteIndicesCount, webgl.UNSIGNED_SHORT, 0)
}

func (batch *DrawCallBatch) AddSprite() *Sprite {
	if batch.spritesCount >= batch.maxSpritesCount {
		return nil
	}

	batch.spritesCount++
	batch.lastUsedIdx++

	batch.verticesAreModified = true

	return NewSprite(batch, batch.texture, batch.lastUsedIdx)
}
