package ui

import (
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
	OnScreenReady func(s screen.Screen)

	w    screen.Window
	tx   chan screen.Texture
	done chan struct{}

	sz  size.Event
	pos image.Rectangle
	fig tshape
}

type tshape struct {
	col    color.RGBA
	posX   int
	width  int
	posY   int
	height int
}

func (fig *tshape) Draw(w screen.Window) {
	dr1 := image.Rectangle{
		Min: image.Point{
			X: fig.posX - fig.width/2,
			Y: fig.posY - fig.height/4,
		},
		Max: image.Point{
			X: fig.posX + fig.width/2,
			Y: fig.posY,
		},
	}
	dr2 := image.Rectangle{
		Min: image.Point{
			X: fig.posX - fig.height/8,
			Y: fig.posY,
		},
		Max: image.Point{
			X: fig.posX + fig.height/8,
			Y: fig.posY + 3*fig.height/4,
		},
	}

	w.Fill(dr1, fig.col, draw.Src)
	w.Fill(dr2, fig.col, draw.Src)
}

var colRed = color.RGBA{255, 0, 0, 255}
var defaultTShape = tshape{colRed, 0, 400, 0, 400}

func (pw *Visualizer) Main() {
	pw.tx = make(chan screen.Texture)
	pw.done = make(chan struct{})

	pw.sz.HeightPx = defaultHeight
	pw.sz.WidthPx = defaultWidth

	driver.Main(pw.run)
}

func (pw *Visualizer) Update(t screen.Texture) {
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
		pw.OnScreenReady(s)
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

	var t screen.Texture

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

func (pw *Visualizer) handleEvent(e any, t screen.Texture) {

	switch e := e.(type) {

	case lifecycle.Event:
		if e.From == lifecycle.StageDead && e.To == lifecycle.StageAlive {
			f := defaultTShape
			f.posX = 400
			f.posY = 400
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
				f := defaultTShape
				f.posX = int(e.X)
				f.posY = int(e.Y)
				pw.fig = f
			}

			pw.w.Send(paint.Event{})
		}

	case paint.Event:
		if t == nil {
			pw.drawDefaultUI()
		} else {
			// Використання текстури отриманої через виклик Update.
			pw.w.Scale(pw.sz.Bounds(), t, t.Bounds(), draw.Src, nil)
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
		pw.w.Fill(br, color.RGBA{200, 0, 0, 255}, draw.Src)
	}
}
