package lang

import (
	"architecture-lab-3/primitives"
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"architecture-lab-3/painter"
)

type parserFunc func(args []float64) (painter.Operation, error)

// Parser уміє прочитати дані з вхідного io.Reader та повернути список операцій представлені вхідним скриптом.
type Parser struct {
	parsingSet map[string]parserFunc
}

func NewParserDefault() Parser {
	p := Parser{parsingSet: commands}
	return p
}

var commands = map[string]parserFunc{
	"white": func(args []float64) (painter.Operation, error) {
		if len(args) != 0 {
			return nil, errors.New("invalid arguments size")
		}
		return painter.OperationFunc(painter.WhiteFill), nil
	},
	"green": func(args []float64) (painter.Operation, error) {
		if len(args) != 0 {
			return nil, errors.New("invalid arguments size")
		}
		return painter.OperationFunc(painter.GreenFill), nil
	},
	"update": func(args []float64) (painter.Operation, error) {
		if len(args) != 0 {
			return nil, errors.New("invalid arguments size")
		}
		return painter.UpdateOp, nil
	},
	"bgrect": func(args []float64) (painter.Operation, error) {
		if len(args) != 4 {
			return nil, errors.New("invalid arguments size")
		}
		rect := primitives.RectF{X1: args[0], X2: args[1], Y1: args[2], Y2: args[3]}
		return painter.OperationFunc(func(ts primitives.TextureStateI) {
			painter.BgRect(ts, rect)
		}), nil
	},
	"figure": func(args []float64) (painter.Operation, error) {
		if len(args) != 2 {
			return nil, errors.New("invalid arguments size")
		}
		point := primitives.PointF{X: args[0], Y: args[1]}
		return painter.OperationFunc(func(ts primitives.TextureStateI) {
			painter.Figure(ts, point)
		}), nil
	},
	"move": func(args []float64) (painter.Operation, error) {
		if len(args) != 2 {
			return nil, errors.New("invalid arguments size")
		}
		point := primitives.PointF{X: args[0], Y: args[1]}
		return painter.OperationFunc(func(ts primitives.TextureStateI) {
			painter.Move(ts, point)
		}), nil
	},
	"reset": func(args []float64) (painter.Operation, error) {
		if len(args) != 0 {
			return nil, errors.New("invalid arguments size")
		}
		return painter.OperationFunc(painter.Reset), nil
	},
}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	var res []painter.Operation

	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		commandLine := scanner.Text()
		op, err := p.parseString(commandLine) // parse the line to get Operation

		if err != nil {
			return nil, err
		}

		res = append(res, op)
	}

	return res, nil
}

func (p *Parser) parseString(inStr string) (painter.Operation, error) {
	commandParams := strings.Fields(inStr)

	cmd, argStrings, err1 := p.getOperationAndArguments(commandParams)
	if err1 != nil {
		return nil, fmt.Errorf("getting operation error: %v", err1)
	}

	args, err2 := p.convertArguments(argStrings)
	if err2 != nil {
		return nil, fmt.Errorf("getting arguments error: %v", err2)
	}

	op, err3 := cmd(args)
	if err3 != nil {
		return nil, fmt.Errorf("creating operation error: %v", err3)
	}

	return op, nil
}

func (p *Parser) getOperationAndArguments(commandParams []string) (parserFunc, []string, error) {
	if len(commandParams) == 0 {
		return nil, nil, errors.New("empty command")
	}

	cmd, ok := p.parsingSet[commandParams[0]]
	if !ok {
		return nil, nil, errors.New("unknown command")
	}

	return cmd, commandParams[1:], nil
}

func (p *Parser) convertArguments(argStrings []string) ([]float64, error) {
	var args []float64
	for _, argStr := range argStrings {
		arg, err := strconv.ParseFloat(argStr, 64)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	return args, nil
}
