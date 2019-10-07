// +build js

package render

import (
	"github.com/nuberu/webgl"
	"github.com/nuberu/webgl/types"
	"strconv"
	"syscall/js"
	"time"
)

type (
	RenderFunc func(dt float32)
	ClickFunc  func(pos Point)

	WebGLRender struct {
		renderFrameCb js.Func

		jsGlobal  js.Value
		jsCanvas  js.Value
		canvasDPR float32
		glCtx     *webgl.RenderingContext

		textRenderer *TextRenderer

		canvasWidth  float32
		canvasHeight float32

		renderFunc     RenderFunc
		lastRenderTime time.Time

		clickFunc ClickFunc
	}
)

const (
	// 4 номера вершин прямоугольника (двух треугольников)
	SpriteRectIndices = 4
	// 6 номеров вершин прямоугольника (двух треугольников)
	SpriteIndicesCount = 6
	// число координат у каждой вершины: x, y
	SpritePosVertices = 2
	// SpritePosVertices координат каждой из вершин
	VerticesSpriteSize = SpritePosVertices * SpriteRectIndices

	// число координат у каждой вершины текстуры: u, v
	SpriteTextureVertices = 2
	// Число значений на описание всех текстурных координат прямоугольника
	VerticesTextureSpriteSize = SpriteTextureVertices * SpriteRectIndices

	// число координат у каждого цвета: r, g, b, a
	SpriteColorsVertices = 4
	// Число значений на описание всех RGBA цветов прямоугольника
	VerticesColorsSpriteSize = SpriteColorsVertices * SpriteRectIndices
)

func NewWebGLRender(canvasId, canvasTextId string) (*WebGLRender, error) {
	rend := &WebGLRender{
		jsGlobal: js.Global(),
	}

	rend.jsCanvas = rend.jsGlobal.Get(`document`).Call(`getElementById`, canvasId)

	var err error
	rend.glCtx, err = webgl.FromCanvas(rend.jsCanvas)
	if err != nil {
		return nil, err
	}

	rend.textRenderer = NewTextRenderer(canvasTextId)

	rend.resizeCanvas(rend.jsCanvas, rend.textRenderer.jsCanvas)
	rend.textRenderer.Resize()

	rend.setupInput()

	rend.setupRenderFunc()

	return rend, nil
}

func (rend *WebGLRender) resizeCanvas(canvases ...js.Value) {
	document := rend.jsGlobal.Get(`document`)
	body_ := document.Get(`body`)
	html_ := document.Get(`documentElement`)

	max := func(vals ...int) int {
		m := 0
		for _, val := range vals {
			if m < val {
				m = val
			}
		}
		return m
	}

	width := max(
		body_.Get(`scrollWidth`).Int(),
		body_.Get(`offsetWidth`).Int(),
		html_.Get(`clientWidth`).Int(),
		html_.Get(`scrollWidth`).Int(),
		html_.Get(`offsetWidth`).Int(),
	)

	height := max(
		body_.Get(`scrollHeight`).Int(),
		body_.Get(`offsetHeight`).Int(),
		html_.Get(`clientHeight`).Int(),
		html_.Get(`scrollHeight`).Int(),
		html_.Get(`offsetHeight`).Int(),
	)

	// HiDPI
	dpr := float32(rend.jsGlobal.Get(`window`).Get(`devicePixelRatio`).Float())
	if dpr == 0 {
		dpr = 1
	}
	rend.canvasDPR = dpr

	rend.canvasWidth = float32(width) * dpr
	rend.canvasHeight = float32(height) * dpr

	for _, canvas := range canvases {
		canvas.Set(`width`, float64(rend.canvasWidth))
		canvas.Set(`height`, float64(rend.canvasHeight))

		if style := canvas.Get(`style`); true {
			style.Set(`width`, strconv.Itoa(width)+`px`)
			style.Set(`height`, strconv.Itoa(height)+`px`)
		}
	}
}

func (rend *WebGLRender) setupInput() {
	rend.jsCanvas.Call(`addEventListener`, `click`, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if rend.clickFunc != nil {
			event := args[0]
			x := float32(event.Get(`clientX`).Float()) * rend.canvasDPR
			y := float32(event.Get(`clientY`).Float()) * rend.canvasDPR
			rend.clickFunc(Point{X: x, Y: y})
		}
		return nil
	}))
}

func (rend *WebGLRender) setupRenderFunc() {
	requestAnimationFrame := rend.jsGlobal.Get(`requestAnimationFrame`)

	rend.lastRenderTime = time.Now()

	rend.renderFrameCb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if rend.renderFunc != nil {
			now := time.Now()
			dt := now.Sub(rend.lastRenderTime).Seconds()
			rend.renderFunc(float32(dt))
			rend.lastRenderTime = now
		}
		requestAnimationFrame.Invoke(rend.renderFrameCb)
		return nil
	})

	requestAnimationFrame.Invoke(rend.renderFrameCb)
}

func (rend *WebGLRender) SetRenderFunc(renderFunc RenderFunc) {
	rend.renderFunc = renderFunc
}

func (rend *WebGLRender) SetClickFunc(clickFunc ClickFunc) {
	rend.clickFunc = clickFunc
}

func (rend *WebGLRender) Width() float32 {
	return rend.canvasWidth
}

func (rend *WebGLRender) Height() float32 {
	return rend.canvasHeight
}

func (rend *WebGLRender) Size() (width float32, height float32) {
	return rend.canvasWidth, rend.canvasHeight
}

func (rend *WebGLRender) GetGlCtx() *webgl.RenderingContext {
	return rend.glCtx
}

func (rend *WebGLRender) GetShaderInfoLog(shader *types.Shader) string {
	return rend.glCtx.GetShaderInfoLog(shader)
}

func (rend *WebGLRender) CompileVertexShader(source string) *types.Shader {
	gl := rend.glCtx
	shader := gl.CreateVertexShader()
	gl.ShaderSource(shader, source)
	gl.CompileShader(shader)
	if err := rend.GetShaderInfoLog(shader); err != `` {
		println(`Compile vertex shader error: ` + err + "\nShader code:\n" + source)
		gl.DeleteShader(shader)
		return nil
	}
	return shader
}

func (rend *WebGLRender) CompileFragmentShader(source string) *types.Shader {
	gl := rend.glCtx
	shader := gl.CreateFragmentShader()
	gl.ShaderSource(shader, source)
	gl.CompileShader(shader)
	return shader
}

func (rend *WebGLRender) CreateShaderProgram(vertexShader, fragmentShader *types.Shader) *types.Program {
	gl := rend.glCtx
	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)
	gl.UseProgram(program)
	return program
}

func (rend *WebGLRender) LoadShaderProgram(vertexShader, fragmentShader string) *types.Program {
	return rend.CreateShaderProgram(
		rend.CompileVertexShader(vertexShader),
		rend.CompileFragmentShader(fragmentShader),
	)
}

func (rend *WebGLRender) CreateArrayBuffer(srcData []float32) *types.Buffer {
	gl := rend.glCtx
	idx := gl.CreateBuffer()
	rend.UpdateArrayBuffer(idx, srcData)
	return idx
}

func (rend *WebGLRender) UpdateArrayBuffer(bufferIdx *types.Buffer, srcData []float32) {
	gl := rend.glCtx
	gl.BindBuffer(webgl.ARRAY_BUFFER, bufferIdx)
	gl.BufferData(webgl.ARRAY_BUFFER, srcData, webgl.STATIC_DRAW)
	gl.BindBuffer(webgl.ARRAY_BUFFER, nil)
}

func (rend *WebGLRender) CreateElementArrayBufferUI16(srcData []uint16) *types.Buffer {
	gl := rend.glCtx
	idx := gl.CreateBuffer()
	gl.BindBuffer(webgl.ELEMENT_ARRAY_BUFFER, idx)
	gl.BufferDataUI16(webgl.ELEMENT_ARRAY_BUFFER, srcData, webgl.STATIC_DRAW)
	gl.BindBuffer(webgl.ELEMENT_ARRAY_BUFFER, nil)
	return idx
}

func (rend *WebGLRender) setupVertexBufferAttrib(shaderProgram *types.Program, name string,
	elemSize int, buf *types.Buffer,
) int {
	gl := rend.glCtx
	attribIdx := gl.GetAttribLocation(shaderProgram, name)
	if attribIdx >= 0 {
		gl.BindBuffer(webgl.ARRAY_BUFFER, buf)
		gl.VertexAttribPointer(attribIdx, elemSize, webgl.FLOAT, false, 0, 0)
		gl.EnableVertexAttribArray(attribIdx)
		gl.BindBuffer(webgl.ARRAY_BUFFER, nil)
	}
	return attribIdx
}

func (rend *WebGLRender) SetupPosVertexBufferAttrib(shaderProgram *types.Program, buf *types.Buffer) int {
	return rend.setupVertexBufferAttrib(shaderProgram, `in_Position`, SpritePosVertices, buf)
}

func (rend *WebGLRender) SetupTextureVertexBufferAttrib(shaderProgram *types.Program, buf *types.Buffer) int {
	return rend.setupVertexBufferAttrib(shaderProgram, `in_TexCoords`, SpriteTextureVertices, buf)
}

func (rend *WebGLRender) SetupColorsVertexBufferAttrib(shaderProgram *types.Program, buf *types.Buffer) int {
	return rend.setupVertexBufferAttrib(shaderProgram, `in_Color`, SpriteColorsVertices, buf)
}

func (rend *WebGLRender) SetupPVUniform(shaderProgram *types.Program) *types.UniformLocation {
	gl := rend.glCtx
	loc := gl.GetUniformLocation(shaderProgram, `PV`)
	if loc != nil {
		var pv Matrix
		pv.Identity()
		w, h := rend.Size()
		pv.Translate(-1, 1)
		pv.Scale(1/(w/2), 1/(-h/2))
		gl.UniformMatrix3fv(loc, false, pv.Val[:])
	}
	return loc
}

func (rend *WebGLRender) SetupSamplerUniform(shaderProgram *types.Program, texIdx int) *types.UniformLocation {
	gl := rend.glCtx
	loc := gl.GetUniformLocation(shaderProgram, `uSampler`)
	if loc != nil {
		// gl.Uniform1i(uniformLoc, int32(tex.texUnit - gl.TEXTURE0))
		gl.Uniform1i(loc, texIdx)
	}
	return loc
}

//func (rend *WebGLRender) LoadTextureSync(url_ string) *Texture {
//	return NewTextureSync(rend.glCtx, url_)
//}

func (rend *WebGLRender) LoadTexture(url_ string, cb TextureLoadCallback) {
	NewTextureAsync(rend.glCtx, url_, cb)
}

func (rend *WebGLRender) Clear() {
	gl := rend.glCtx
	gl.ClearColor(1, 1, 1, 1)
	gl.Clear(uint32(webgl.COLOR_BUFFER_BIT))

	rend.textRenderer.Draw()
}

func (rend *WebGLRender) AddText(text string, size int, pos Point) *Text2D {
	return rend.textRenderer.AddText(text, size, pos)
}

func (rend *WebGLRender) DeleteText(text *Text2D) {
	rend.textRenderer.DeleteText(text)
}
