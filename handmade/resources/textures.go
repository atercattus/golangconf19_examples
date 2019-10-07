package resources

import "github.com/atercattus/golangconf19_examples/handmade/render"

const (
	imgFilenameGoRaces    = `./images/go_races.png`
	imgFilenameRunBtn     = `./images/shadedDark25.png`
	imgFilenameGopher     = `./images/gopher_no_paw_scooker.png`
	imgFilenameRoad       = `./images/tunnel_road.jpg`
	imgFilenameGopherPaw  = `./images/paw.png`
	imgFilenameBtn        = `./images/btn.png`
	imgFilenameBkg        = `./images/backgroundEmpty.png`
	imgFilenameFinishLine = `./images/finish_line.png`
)

var (
	TextureTitle      *render.Texture
	TextureRunBtn     *render.Texture
	TextureGopher     *render.Texture
	TextureGopherPaw  *render.Texture
	TextureRoad       *render.Texture
	TextureBtn        *render.Texture
	TextureBkg        *render.Texture
	TextureFinishLine *render.Texture
)

func TexturesPreload(renderer *render.WebGLRender, done func()) {
	renderer.LoadTexture(imgFilenameGoRaces, func(texture *render.Texture) {
		TextureTitle = texture

		renderer.LoadTexture(imgFilenameRunBtn, func(texture *render.Texture) {
			TextureRunBtn = texture

			renderer.LoadTexture(imgFilenameGopher, func(texture *render.Texture) {
				TextureGopher = texture

				renderer.LoadTexture(imgFilenameGopherPaw, func(texture *render.Texture) {
					TextureGopherPaw = texture

					renderer.LoadTexture(imgFilenameRoad, func(texture *render.Texture) {
						TextureRoad = texture

						renderer.LoadTexture(imgFilenameBtn, func(texture *render.Texture) {
							TextureBtn = texture

							renderer.LoadTexture(imgFilenameBkg, func(texture *render.Texture) {
								TextureBkg = texture

								renderer.LoadTexture(imgFilenameFinishLine, func(texture *render.Texture) {
									TextureFinishLine = texture

									done()
								})
							})
						})
					})
				})
			})
		})
	})
}
