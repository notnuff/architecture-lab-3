package main

import (
	"architecture-lab-3/painter"
	"architecture-lab-3/painter/lang"
	"architecture-lab-3/ui"
	"net/http"
)

func main() {
	var (
		pv ui.Visualizer // Візуалізатор створює вікно та малює у ньому.

		// Потрібні для частини 2.
		opLoop painter.Loop              // Цикл обробки команд.
		parser = lang.NewParserDefault() // Парсер команд.
	)

	pv.Debug = true
	pv.Title = "Simple painter"

	pv.OnScreenReady = opLoop.Start
	opLoop.Receiver = &pv

	go func() {
		http.Handle("/", lang.HttpHandler(&opLoop, &parser))
		_ = http.ListenAndServe("localhost:17000", nil)
	}()

	//go func() {
	//	time.Sleep(time.Second)
	//
	//	opLoop.Post(painter.OperationFunc(painter.GreenFill))
	//	opLoop.Post(painter.UpdateOp)
	//
	//	time.Sleep(time.Second)
	//	opLoop.Post(painter.OperationFunc(func(ts primitives.TextureStateI) {
	//		coords := primitives.RectF{
	//			X1: 0.1,
	//			X2: 0.2,
	//			Y1: 0.3,
	//			Y2: 0.4,
	//		}
	//		painter.BgRect(ts, coords)
	//	}))
	//	opLoop.Post(painter.UpdateOp)
	//
	//	time.Sleep(time.Second)
	//	opLoop.Post(painter.OperationFunc(func(ts primitives.TextureStateI) {
	//		coords := primitives.RectF{
	//			X1: 0.3,
	//			X2: 0.1,
	//			Y1: 0.5,
	//			Y2: 0.7,
	//		}
	//		painter.BgRect(ts, coords)
	//	}))
	//	opLoop.Post(painter.UpdateOp)
	//
	//	time.Sleep(time.Second)
	//	opLoop.Post(painter.OperationFunc(func(ts primitives.TextureStateI) {
	//		coords := primitives.PointF{
	//			X: 0.6,
	//			Y: 0.6,
	//		}
	//		painter.Figure(ts, coords)
	//	}))
	//	opLoop.Post(painter.UpdateOp)
	//
	//	time.Sleep(time.Second)
	//	opLoop.Post(painter.OperationFunc(func(ts primitives.TextureStateI) {
	//		coords := primitives.PointF{
	//			X: 0.2,
	//			Y: 0.3,
	//		}
	//		painter.Figure(ts, coords)
	//	}))
	//	opLoop.Post(painter.UpdateOp)
	//
	//	time.Sleep(time.Second)
	//	opLoop.Post(painter.OperationFunc(func(ts primitives.TextureStateI) {
	//		coords := primitives.PointF{
	//			X: 0.2,
	//			Y: 0.3,
	//		}
	//		painter.Move(ts, coords)
	//	}))
	//	opLoop.Post(painter.OperationFunc(func(ts primitives.TextureStateI) {
	//		coords := primitives.PointF{
	//			X: -0.3,
	//			Y: -0.4,
	//		}
	//		painter.Move(ts, coords)
	//	}))
	//	opLoop.Post(painter.OperationFunc(painter.WhiteFill))
	//	opLoop.Post(painter.UpdateOp)
	//
	//	time.Sleep(time.Second)
	//	opLoop.Post(painter.OperationFunc(func(ts primitives.TextureStateI) {
	//		painter.Reset(ts)
	//	}))
	//	opLoop.Post(painter.UpdateOp)
	//}()

	pv.Main()
	opLoop.StopAndWait()

}
