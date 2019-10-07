package common

import (
	"github.com/EngoEngine/engo/common"
	"image"
	"image/color"
)

func NewColorTexture(clr color.Color, width float32, height float32) *common.Texture {
	img := image.NewUniform(clr)
	nrgba := common.ImageToNRGBA(img, int(width), int(height))
	imageObj := common.NewImageObject(nrgba)
	texture := common.NewTextureSingle(imageObj)
	return &texture
}
