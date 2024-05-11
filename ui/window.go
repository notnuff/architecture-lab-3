package ui

import (
	"golang.org/x/exp/shiny/driver"
	"image"
	"image/color"
	"log"
	"time"

	"golang.org/x/exp/shiny/imageutil"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
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
}

type changeColor struct {
	C color.RGBA
}

func animate(w screen.Window, duration time.Duration) {
	//timeStart := time.Now()
	var gray uint8 = 50
	for {
		time.Sleep(1000)
		//timePassed := time.Now().Sub(timeStart)
		//timePassed %= duration

		//gray = uint8(timePassed/duration) * 255
		col := changeColor{C: color.RGBA{gray, gray, gray, gray}}
		w.Send(col)
	}
}

func (pw *Visualizer) Main() {
	pw.tx = make(chan screen.Texture)
	pw.done = make(chan struct{})
	pw.pos.Max.X = 0
	pw.pos.Max.Y = 0

	pw.sz.HeightPx = defaultHeight
	pw.sz.WidthPx = defaultWidth

	//driver.Main(pw.run)

	driver.Main(func(s screen.Screen) {
		w, err := s.NewWindow(&screen.NewWindowOptions{
			Width:  800,
			Height: 800,
			Title:  "Gamedev GO on",
		})
		if err != nil {
			log.Fatal(err)
		}

		defer func() {
			w.Release()
			log.Printf("Wow, u closed da window")
		}()

		var (
			sz          size.Event
			m           mouse.Button
			bgc               = color.RGBA{0, 0, 0, 0}
			gray        uint8 = 0
			isAnimating       = false
		)

		for {
			e := w.NextEvent()
			switch e := e.(type) {
			case size.Event:
				if !isAnimating {
					go animate(w, 5*time.Second)
					isAnimating = true
				}

				sz = e
				log.Printf("%i", e.HeightPx)
			case paint.Event:
				w.Fill(sz.Bounds(), bgc, screen.Src)
				w.Publish()
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					w.Release()
					log.Printf("No, fuck u, i`m not closing, %s", e.String())
				}
			case mouse.Event:
				if e.Direction == mouse.DirPress {
					m = e.Button
				}
				if e.Direction == mouse.DirRelease {
					m = mouse.ButtonNone
				}

				if m == mouse.ButtonRight {
					bgc = color.RGBA{uint8(e.X), uint8(e.Y), uint8(e.X + e.Y), 100}
					//w.Fill(sz.Bounds(), bgc, screen.Src)
					w.Send(paint.Event{})
				}
			case key.Event:
				repaint := false
				if e.Code == key.CodeUpArrow {
					gray += 10
					repaint = true
				}
				if e.Code == key.CodeDownArrow {
					gray -= 10
					repaint = true
				}
				if repaint {
					bgc = color.RGBA{gray, gray, gray, gray}
					w.Send(changeColor{})
				}
			case changeColor:
				bgc = e.C
				w.Send(paint.Event{})
			}

		}
	})
}

func (pw *Visualizer) Update(t screen.Texture) {
	pw.tx <- t
}

func (pw *Visualizer) DrawT(posX, posY float32) {
	pw.w.Fill(pw.sz.Bounds(), color.RGBA{uint8(posX), uint8(posY), 6, 10}, screen.Src)

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

	case size.Event: // Оновлення даних про розмір вікна.
		pw.sz = e

	case error:
		log.Printf("ERROR: %s", e)

	case mouse.Event:
		if t == nil {
			if e.Button == mouse.ButtonRight {
				pw.DrawT(e.X, e.Y)
			}
			// TODO: Реалізувати реакцію на натискання кнопки миші.
		}

	case paint.Event:
		// Малювання контенту вікна.
		if t == nil {
			pw.drawDefaultUI()
		} else {
			// Використання текстури отриманої через виклик Update.
			pw.w.Scale(pw.sz.Bounds(), t, t.Bounds(), draw.Src, nil)
		}
		pw.w.Publish()
	}
}

func (pw *Visualizer) drawDefaultUI() {
	pw.w.Fill(pw.sz.Bounds(), color.Black, draw.Src) // Фон.

	// TODO: Змінити колір фону та додати відображення фігури у вашому варіанті.

	// Малювання білої рамки.
	for _, br := range imageutil.Border(pw.sz.Bounds(), 10) {
		pw.w.Fill(br, color.White, draw.Src)
	}
}
