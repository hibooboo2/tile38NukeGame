package main

import (
	"image"
	"os"

	"github.com/hibooboo2/tile38NukeGame/game"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/style"
	"github.com/tidwall/pinhole"
)

var count int
var opts = pinhole.DefaultImageOptions
var n = 60

func main() {
	name, _ := os.Hostname()
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	player = game.NewCharacter(name)

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

var player *game.Character

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
