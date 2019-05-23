package main

import (
	"fmt"
	"image/color"
	"math"
	"sync"
	"time"

	"github.com/aarzilli/nucular"
	"github.com/tidwall/pinhole"
)

var count int
var p = pinhole.New()
var lock = sync.Mutex{}
var opts = pinhole.DefaultImageOptions
var n = 60

func main2() {
	// wnd := nucular.NewMasterWindow(0, "Counter", updatefn)
	// s := style.FromTheme(style.DarkTheme, 2.0)
	// s.NormalWindow.MinSize = image.Point{1920, 1080}
	// go func() {
	// 	for {
	// 		select {
	// 		case <-t:
	// 			next()
	// 		}
	// 	}
	// }()
	// wnd.SetStyle(s)

	// wnd.Main()
	opts.LineWidth = 4
	opts.Scale = 4
}

func updatefn(w *nucular.Window) {

	w.Row(1080).Static(1920)
	height := 1080 / 2
	width := 1920 / 2
	lock.Lock()
	w.Image(p.Image(width, height, opts))
	lock.Unlock()
}

var t = time.Tick(time.Millisecond * 13)

func next() {
	lock.Lock()
	count++
	i := count % n
	fmt.Printf("frame %d/%d\n", i, n)
	p = pinhole.New()
	p.Begin()
	p.DrawCube(-0.2, -0.2, -0.2, 0.2, 0.2, 0.2)
	p.Rotate(0, math.Pi*2/(float64(n)/float64(i)), 0)
	p.Colorize(color.RGBA{255, 0, 0, 255})
	p.End()

	p.Begin()
	p.DrawCircle(0, 0, 0, 0.2)
	p.Rotate(math.Pi*2/(float64(n)/float64(i)), math.Pi*4/(float64(n)/float64(i)), 0)
	p.End()

	p.Begin()
	p.DrawCircle(0, 0, 0, 0.2)
	p.Rotate(-math.Pi*2/(float64(n)/float64(i)), math.Pi*4/(float64(n)/float64(i)), 0)
	p.End()

	p.Scale(1.75, 1.75, 1.75)
	lock.Unlock()
}
