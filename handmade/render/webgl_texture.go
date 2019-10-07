package render

import (
	"github.com/nuberu/webgl"
	"github.com/nuberu/webgl/types"
	"syscall/js"
)

type (
	Texture struct {
		tex *types.Texture
		gl  *webgl.RenderingContext

		size Point
	}

	TextureLoadCallback func(*Texture)
)

func isPowerOfTwo(val int) bool {
	return (val & (val - 1)) == 0
}

// https://developer.mozilla.org/ru/docs/Web/API/WebGL_API/Tutorial/Using_textures_in_WebGL
func NewTextureAsync(gl *webgl.RenderingContext, url_ string, cb TextureLoadCallback) {
	tex := gl.CreateTexture()
	gl.BindTexture(webgl.TEXTURE_2D, tex)

	// чтобы текстуру уже как-то можно было использовать (но надо ли?)
	gl.TexImage2Db(webgl.TEXTURE_2D, 0, webgl.RGBA, 1, 1, 0, webgl.RGBA, []byte{0, 100, 255, 255})

	img := js.Global().Get(`Image`).New()
	img.Set(`onload`, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		width := img.Get(`width`).Int()
		height := img.Get(`height`).Int()

		gl.BindTexture(webgl.TEXTURE_2D, tex)

		gl.TexImage2DHtmlElement(webgl.TEXTURE_2D, 0, webgl.RGBA, webgl.RGBA, webgl.UNSIGNED_BYTE, img)

		gl.TexParameterWrapS(webgl.TEXTURE_2D, webgl.CLAMP_TO_EDGE)
		gl.TexParameterWrapT(webgl.TEXTURE_2D, webgl.CLAMP_TO_EDGE)
		gl.TexParameterMinFilter(webgl.TEXTURE_2D, webgl.NEAREST)
		gl.TexParameterMagFilter(webgl.TEXTURE_2D, webgl.NEAREST)

		if isPowerOfTwo(width) && isPowerOfTwo(height) {
			gl.GenerateMipmap(webgl.TEXTURE_2D)
		}

		cb(&Texture{
			tex:  tex,
			gl:   gl,
			size: Point{X: float32(width), Y: float32(height)},
		})

		return nil
	}))
	img.Set(`onerror`, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		gl.DeleteTexture(tex)
		cb(nil)

		return nil
	}))
	img.Set(`src`, url_)
}

func (t *Texture) Bind() {
	t.gl.BindTexture(webgl.TEXTURE_2D, t.tex)
}

func (t *Texture) Size() Point {
	return t.size
}
