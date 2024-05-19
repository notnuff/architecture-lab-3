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

	//pv.Debug = true
	pv.Title = "Simple painter"

	pv.OnScreenReady = opLoop.Start
	opLoop.Receiver = &pv

	go func() {
		http.Handle("/", lang.HttpHandler(&opLoop, &parser))
		_ = http.ListenAndServe("localhost:17000", nil)
	}()
	pv.Main()
	opLoop.StopAndWait()

}
