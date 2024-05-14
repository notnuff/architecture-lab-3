package painter

import (
	"architecture-lab-3/primitives"
	"image/color"
	"reflect"
	"testing"
)

func TestLoop_Post(t *testing.T) {
	var (
		l  Loop
		tr testReceiver
	)
	l.Receiver = &tr

	var testOps []string

	l.Start()
	l.Post(logOp(t, "do white fill", WhiteFill))
	l.Post(logOp(t, "do green fill", GreenFill))
	l.Post(UpdateOp)

	for i := 0; i < 3; i++ {
		go l.Post(logOp(t, "do green fill", GreenFill))
	}

	l.Post(OperationFunc(func(state primitives.TextureState) {
		testOps = append(testOps, "op 1")

		l.Post(OperationFunc(func(state primitives.TextureState) {
			testOps = append(testOps, "op 3")
		}))
	}))

	l.Post(OperationFunc(func(state primitives.TextureState) {
		testOps = append(testOps, "op 2")
	}))

	l.StopAndWait()

	textureState := tr.ts

	if textureState.BgCol != color.White {
		t.Error("First color is not white:", textureState.BgCol)
	}

	if !reflect.DeepEqual(testOps, []string{"op 1", "op 2", "op 3"}) {
		t.Error("Bad order:", testOps)
	}
}

func logOp(t *testing.T, msg string, op OperationFunc) OperationFunc {
	return func(ts primitives.TextureState) {
		t.Log(msg)
		op(ts)
	}
}

type testReceiver struct {
	ts primitives.TextureState
}

func (tr *testReceiver) Update(t primitives.TextureState) {
	tr.ts = t
}
