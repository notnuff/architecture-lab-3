package primitives

import (
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"image"
	"image/color"
)

type TShape struct {
	Col    color.RGBA
	PosX   int
	PosY   int
	Width  int
	Height int
}

func (fig *TShape) Draw(w screen.Window) {
	dr1 := image.Rectangle{
		Min: image.Point{
			X: fig.PosX - fig.Width/2,
			Y: fig.PosY - fig.Height/4,
		},
		Max: image.Point{
			X: fig.PosX + fig.Width/2,
			Y: fig.PosY,
		},
	}
	dr2 := image.Rectangle{
		Min: image.Point{
			X: fig.PosX - fig.Height/8,
			Y: fig.PosY,
		},
		Max: image.Point{
			X: fig.PosX + fig.Height/8,
			Y: fig.PosY + 3*fig.Height/4,
		},
	}

	w.Fill(dr1, fig.Col, draw.Src)
	w.Fill(dr2, fig.Col, draw.Src)
}

func NewTShape(px, py int) TShape {
	var colRed = color.RGBA{R: 255, A: 255}
	return TShape{colRed, px, py, 400, 400}
}
