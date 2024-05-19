package lang

import (
	"architecture-lab-3/primitives"
	"image/color"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	parser := NewParserDefault()
	{
		command := "white"
		var reader io.Reader = strings.NewReader(command)
		ops, err := parser.Parse(reader)
		if err != nil {
			t.Fatal(err)
		}
		if len(ops) != 1 {
			t.Error("unexpected operations length: ", len(ops))
		}
		op := ops[0]
		testTextureState := new(primitives.TextureState)
		op.Do(testTextureState)
		if testTextureState.BgCol != color.White {
			t.Error("unexpected bgcol value: ", testTextureState.BgCol)
		}
	}
	{
		command := "green"
		var reader io.Reader = strings.NewReader(command)
		ops, err := parser.Parse(reader)
		if err != nil {
			t.Fatal(err)
		}
		if len(ops) != 1 {
			t.Error("unexpected operations length: ", len(ops))
		}
		op := ops[0]
		testTextureState := new(primitives.TextureState)
		op.Do(testTextureState)
		if testTextureState.BgCol != color.Color(color.RGBA{G: 255, A: 255}) {
			t.Error("unexpected bgcol value: ", testTextureState.BgCol)
		}
	}
	{
		command := "bgrect 0.1 0.4 0.2 0.6"
		var reader io.Reader = strings.NewReader(command)
		ops, err := parser.Parse(reader)
		if err != nil {
			t.Fatal(err)
		}
		if len(ops) != 1 {
			t.Error("unexpected operations length: ", len(ops))
		}
		op := ops[0]

		testTextureState := new(primitives.TextureState)
		op.Do(testTextureState)
		if *testTextureState.BgRect != (primitives.RectF{X1: 0.1, X2: 0.4, Y1: 0.2, Y2: 0.6}) {
			t.Error("unexpected BgRect:", testTextureState.BgRect)
		}
	}
	{
		command := "figure 0.1 0.4\nfigure 0.2 0.2\nfigure 0.6 0.5\n"
		var readerFig io.Reader = strings.NewReader(command)
		ops, err := parser.Parse(readerFig)
		if err != nil {
			t.Fatal(err)
		}
		if len(ops) != 3 {
			t.Error("unexpected operations length: ", len(ops))
		}
		testTextureState := new(primitives.TextureState)

		for _, op := range ops {
			op.Do(testTextureState)
		}
		expectedFigures := []primitives.PointF{
			{0.1, 0.4},
			{0.2, 0.2},
			{0.6, 0.5},
		}

		for i, figPos := range testTextureState.Figs {
			if !reflect.DeepEqual(*figPos, expectedFigures[i]) {
				t.Error("Wrong sequence of figures positions:", testTextureState.Figs)
			}
		}

		commandMove := "move -0.1 0.1"
		var readerMove io.Reader = strings.NewReader(commandMove)
		opsMove, errMove := parser.Parse(readerMove)
		if errMove != nil {
			t.Fatal(err)
		}
		if len(opsMove) != 1 {
			t.Error("unexpected operations length: ", len(ops))
		}
		opsMove[0].Do(testTextureState)
		for i, figPos := range testTextureState.Figs {
			movedPos := primitives.PointF{expectedFigures[i].X - 0.1, expectedFigures[i].Y + 0.1}
			if !reflect.DeepEqual(*figPos, movedPos) {
				t.Error("Wrong sequence of figures positions:", testTextureState.Figs)
			}
		}
	}
	{
		command := "bgrect 0.1 0.4 0.2 0.6 0.4"
		var reader io.Reader = strings.NewReader(command)
		_, err := parser.Parse(reader)
		if err == nil {
			t.Fatal("parser should have failed with invalid arguments")
		}
	}
	{
		command := "figure 0.1 0.6 0.2"
		var reader io.Reader = strings.NewReader(command)
		_, err := parser.Parse(reader)
		if err == nil {
			t.Fatal("parser should have failed with invalid arguments")
		}
	}
	{
		command := "move 0.1"
		var reader io.Reader = strings.NewReader(command)
		_, err := parser.Parse(reader)
		if err == nil {
			t.Fatal("parser should have failed with invalid arguments")
		}
	}

}

func TestParser_convertArguments(t *testing.T) {
	parser := NewParserDefault()
	tests := []struct {
		name    string
		args    []string
		want    []float64
		wantErr bool
	}{
		{"empty test",
			strings.Fields(""),
			nil,
			false,
		},
		{"one argument conversion",
			strings.Fields("0.1"),
			[]float64{0.1},
			false,
		},
		{"two arguments conversion",
			strings.Fields("0.1 0.7"),
			[]float64{0.1, 0.7},
			false,
		},
		{"negative arguments conversion",
			strings.Fields("-0.1 -0.7 -0.6 -0.456"),
			[]float64{-0.1, -0.7, -0.6, -0.456},
			false,
		},
		{"wrong arguments format conversion",
			strings.Fields("-0.1.2 7f"),
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.convertArguments(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertArguments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertArguments() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_getOperationAndArguments(t *testing.T) {
	parser := NewParserDefault()
	tests := []struct {
		name           string
		command        []string
		wantArgsString []string
		wantErr        bool
	}{
		{"empty test",
			strings.Fields(""),
			nil,
			true,
		},
		{"operation without args test",
			strings.Fields("green"),
			[]string{},
			false,
		},
		{"operation with args test",
			strings.Fields("figure 0.1 0.2 0.3 0.4"),
			strings.Fields("0.1 0.2 0.3 0.4"),
			false,
		},
		{"unknown operation test",
			strings.Fields("blabla"),
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, argumentsString, err := parser.getOperationAndArguments(tt.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("getOperationAndArguments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(argumentsString, tt.wantArgsString) {
				t.Errorf("getOperationAndArguments() argumentsString = %v, want %v", argumentsString, tt.wantArgsString)
			}
		})
	}
}

func TestParser_parseString(t *testing.T) {
	parser := NewParserDefault()
	tests := []struct {
		name             string
		inString         string
		wantTextureState primitives.TextureState
		wantErr          bool
	}{
		{
			"test empty",
			"",
			primitives.TextureState{},
			true,
		},
		{
			"green background test",
			"green",
			primitives.TextureState{
				BgCol:  color.Color(color.RGBA{G: 255, A: 255}),
				BgRect: nil,
				Figs:   nil,
			},
			false,
		},
		{
			"wrong arguments test",
			"green 0.5",
			primitives.TextureState{},
			true,
		},
		{
			"figure test",
			"figure 0.5 0.3",
			primitives.TextureState{
				BgCol:  nil,
				BgRect: nil,
				Figs: []*primitives.PointF{
					{X: 0.5, Y: 0.3},
				},
			},
			false,
		},
		{
			"bgrect test",
			"bgrect 0.5 0.3 0.6 0.6",
			primitives.TextureState{
				BgCol:  nil,
				BgRect: &primitives.RectF{0.5, 0.3, 0.6, 0.6},
				Figs:   nil,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op, err := parser.parseString(tt.inString)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if op == nil {
				return
			}

			testTS := primitives.TextureState{}
			op.Do(&testTS)
			if !reflect.DeepEqual(testTS, tt.wantTextureState) {
				t.Errorf("parseString() textureState = %v, wantTextureState %v", testTS, tt.wantTextureState)
				return
			}
		})
	}
}
