package primitives

import (
	"image/color"
)

type TextureState struct {
	BgCol  color.Color
	BgRect *RectF
	Figs   []*PointF
}

func (t *TextureState) SetBgRect(rect RectF) {
	if t.BgRect == nil {
		t.BgRect = new(RectF)
	}
	*t.BgRect = rect

}
func (t *TextureState) SetBgColor(color color.Color) {
	t.BgCol = color
}
func (t *TextureState) AddBgFig(pos PointF) {
	posP := new(PointF)
	posP.X = pos.X
	posP.Y = pos.Y
	t.Figs = append(t.Figs, posP)
}

func (t *TextureState) GetBgColor() color.Color {
	return t.BgCol
}

func (t *TextureState) GetBgRect() *RectF {
	return t.BgRect
}

func (t *TextureState) GetFigs() []*PointF {
	return t.Figs
}

func (t *TextureState) ResetFigs() {
	clear(t.Figs)
	t.Figs = []*PointF{}
	t.BgRect = nil
}

type TextureStateI interface {
	SetBgRect(RectF)
	SetBgColor(color.Color)
	AddBgFig(PointF)

	GetBgColor() color.Color
	GetBgRect() *RectF
	GetFigs() []*PointF
	ResetFigs()
}
