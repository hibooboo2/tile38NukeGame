// +build nuke

package main

import (
	"image"
	"image/color"
	"log"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/style"
)

func ui() {
	log.Println("NUKE!")
	wnd := nucular.NewMasterWindow(0, "Counter", updatefn)
	s := style.FromTheme(style.DarkTheme, 2.0)
	s.NormalWindow.MinSize = image.Point{1920, 1080}
	wnd.SetStyle(s)

	wnd.Main()
	opts.LineWidth = 4
	opts.Scale = 4
}

func updatefn(w *nucular.Window) {
	handleMoveUsers(w.Input())

	w.Row(200).Static(200)
	img := imgWrap{image.NewRGBA(image.Rect(0, 0, 1000, 1000))}

	player.MiniMap(img.drawBox, 200, 200, 200)
	w.Image(img.RGBA)
}

func handleMoveUsers(w *nucular.Input) {
	for _, key := range w.Keyboard.Keys {
		switch key.Rune {
		case 'w':
			player.MoveRel(0, 1)
		case 'a':
			player.MoveRel(1, 0)
		case 's':
			player.MoveRel(0, -1)
		case 'd':
			player.MoveRel(-1, 0)
		}
	}
}

type imgWrap struct {
	*image.RGBA
}

func (img *imgWrap) drawBox(x, y, size int, c color.RGBA) {
	b := img.Bounds()
	if x > b.Min.X && x < b.Max.X && y > b.Min.Y && y < b.Max.Y {
		for i := -size / 2; i < size/2; i++ {
			for j := -size / 2; j < size/2; j++ {
				img.Set(x+i, y+j, c)
			}
		}
	}
}
