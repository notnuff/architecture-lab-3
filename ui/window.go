package ui

import (
	"architecture-lab-3/primitives"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/imageutil"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"image"
	"image/color"
	"log"
)

const (
	defaultWidth  = 800
	defaultHeight = 800
)

type Visualizer struct {
	Title         string
	Debug         bool
	OnScreenReady func()

	w  screen.Window
	tx chan primitives.TextureStateI

	done chan struct{}

	sz  size.Event
	pos image.Rectangle
	fig primitives.TShape
}

func (pw *Visualizer) Main() {
	pw.tx = make(chan primitives.TextureStateI)
	pw.done = make(chan struct{})

	pw.sz.HeightPx = defaultHeight
	pw.sz.WidthPx = defaultWidth

	driver.Main(pw.run)
}

func (pw *Visualizer) Update(t primitives.TextureStateI) {
	pw.tx <- t
}

func (pw *Visualizer) run(s screen.Screen) { //this function takes control after drivers initialization
	w, err := s.NewWindow(&screen.NewWindowOptions{
		Width:  pw.sz.WidthPx,
		Height: pw.sz.HeightPx,
		Title:  pw.Title,
	})

	if err != nil {
		log.Fatal("Failed to initialize the app window:", err)
	}
	defer func() {
		w.Release()
		close(pw.done)
	}()

	if pw.OnScreenReady != nil {
		pw.OnScreenReady()
	}

	pw.w = w

	events := make(chan any)
	go func() {
		for {
			e := w.NextEvent()
			if pw.Debug {
				log.Printf("new event: %v", e)
			}
			if detectTerminate(e) {
				close(events)
				break
			}
			events <- e
		}
	}()

	var t primitives.TextureStateI

	for {
		select {
		case e, ok := <-events:
			if !ok {
				return
			}
			pw.handleEvent(e, t)

		case t = <-pw.tx:
			w.Send(paint.Event{})
		}
	}
}

func detectTerminate(e any) bool {
	switch e := e.(type) {
	case lifecycle.Event:
		if e.To == lifecycle.StageDead {
			return true // Window destroy initiated.
		}
	case key.Event:
		if e.Code == key.CodeEscape {
			return true // Esc pressed.
		}
	}
	return false
}

func (pw *Visualizer) handleEvent(e any, t primitives.TextureStateI) {

	switch e := e.(type) {

	case lifecycle.Event:
		if e.From == lifecycle.StageDead && e.To == lifecycle.StageAlive {
			f := primitives.NewTShape(400, 400)
			pw.fig = f
			pw.w.Send(paint.Event{})
		}
	case size.Event: // Оновлення даних про розмір вікна.
		pw.sz = e
	case error:
		log.Printf("ERROR: %s", e)

	case mouse.Event:
		if t == nil {
			if e.Direction != mouse.DirPress {
				return
			}

			if e.Button == mouse.ButtonRight {
				f := primitives.NewTShape(int(e.X), int(e.Y))
				pw.fig = f
			}

			pw.w.Send(paint.Event{})
		}

	case paint.Event:

		if t == nil {
			pw.drawDefaultUI()
		} else {
			if bgCol := t.GetBgColor(); bgCol != nil {
				pw.w.Fill(pw.sz.Bounds(), bgCol, draw.Src)
			}
			if bgRect := t.GetBgRect(); bgRect != nil {

				bgRectScaled := image.Rect(
					int(bgRect.X1*float64(pw.sz.WidthPx)),
					int(bgRect.X2*float64(pw.sz.WidthPx)),
					int(bgRect.Y1*float64(pw.sz.HeightPx)),
					int(bgRect.Y2*float64(pw.sz.HeightPx)),
				)
				pw.w.Fill(bgRectScaled, color.Black, draw.Src)
			}

			for _, fig := range t.GetFigs() {

				tshapePosScaled := image.Point{
					X: int(fig.X * float64(pw.sz.WidthPx)),
					Y: int(fig.Y * float64(pw.sz.HeightPx)),
				}

				f := primitives.NewTShape(tshapePosScaled.X, tshapePosScaled.Y)
				f.Draw(pw.w)
			}

			// Використання текстури отриманої через виклик Update.
			//pw.w.Scale(pw.sz.Bounds(), t, t.Bounds(), draw.Src, nil)
		}
		pw.w.Publish()

	}
}

func (pw *Visualizer) drawDefaultBackground() {
	pw.w.Fill(pw.sz.Bounds(), color.White, draw.Src) // Фон.
}

func (pw *Visualizer) drawFigure() {
	pw.fig.Draw(pw.w)
}

func (pw *Visualizer) drawDefaultUI() {
	pw.drawDefaultBackground()
	pw.drawFigure()

	for _, br := range imageutil.Border(pw.sz.Bounds(), 10) {
		pw.w.Fill(br, color.RGBA{R: 200, A: 255}, draw.Src)
	}
}
