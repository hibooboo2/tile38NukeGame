package main

import (
	"image"
	"time"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/style"
	"github.com/tidwall/pinhole"
)

var count int
var opts = pinhole.DefaultImageOptions
var n = 60

func main() {
	go players()
	wnd := nucular.NewMasterWindow(0, "Counter", updatefn)
	s := style.FromTheme(style.DarkTheme, 2.0)
	s.NormalWindow.MinSize = image.Point{1920, 1080}
	wnd.SetStyle(s)

	wnd.Main()
	opts.LineWidth = 4
	opts.Scale = 4
}

func updatefn(w *nucular.Window) {
	height := 1080 / 16
	width := 1920 / 16

	w.Row(height).Static(width)
	height = 1080 / 4
	width = 1920 / 4
	charLock.Lock()
	for _, c := range chars {
		p := c.GetPinHole()
		w.Image(p.Image(width, height, opts))
	}
	charLock.Unlock()
}

var t = time.Tick(time.Millisecond * 13)
