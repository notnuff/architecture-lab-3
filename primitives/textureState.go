package primitives

import (
	"image"
	"image/color"
)

type TextureState struct {
	BgCol  color.Color
	BgRect image.Rectangle
	Figs   []TShape
}
