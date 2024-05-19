package painter

import (
	"architecture-lab-3/primitives"
	"image/color"
)

// Operation змінює вхідну текстуру.
type Operation interface {
	// Do виконує зміну операції, повертаючи true, якщо текстура вважається готовою для відображення.
	Do(ts primitives.TextureStateI) (ready bool)
}

// OperationList групує список операції в одну.
type OperationList []Operation

func (ol OperationList) Do(ts primitives.TextureStateI) (ready bool) {
	for _, o := range ol {
		ready = o.Do(ts) || ready
	}
	return
}

// UpdateOp операція, яка не змінює текстуру, але сигналізує, що текстуру потрібно розглядати як готову.
var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(ts primitives.TextureStateI) bool { return true }

// OperationFunc використовується для перетворення функції оновлення текстури в Operation.
type OperationFunc func(ts primitives.TextureStateI)

func (f OperationFunc) Do(ts primitives.TextureStateI) bool {
	f(ts)

	return false
}

// WhiteFill зафарбовує тестуру у білий колір. Може бути викоистана як Operation через OperationFunc(WhiteFill).
func WhiteFill(ts primitives.TextureStateI) {
	ts.SetBgColor(color.White)
}

// GreenFill зафарбовує тестуру у зелений колір. Може бути викоистана як Operation через OperationFunc(GreenFill).
func GreenFill(ts primitives.TextureStateI) {
	col := color.RGBA{G: 255, A: 255}
	ts.SetBgColor(col)
}

func BgRect(ts primitives.TextureStateI, r primitives.RectF) {
	ts.SetBgRect(r)
}

func Figure(ts primitives.TextureStateI, p primitives.PointF) {
	ts.AddBgFig(p)
}

func Move(ts primitives.TextureStateI, offset primitives.PointF) {
	figsPos := ts.GetFigs()
	for _, pos := range figsPos {
		pos.X += offset.X
		pos.Y += offset.Y
	}
}

func Reset(ts primitives.TextureStateI) {
	ts.SetBgColor(color.Black)
	ts.ResetFigs()
}
