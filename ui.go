// +build nuke

package main

import (
	"image"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/style"
)

func ui() {
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
	img := image.NewRGBA(image.Rect(0, 0, 1000, 1000))
	player.MiniMap(img, 200, 200, 200)
	w.Image(img)
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
