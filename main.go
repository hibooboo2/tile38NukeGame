//build !js,!wasm

package main

import (
	"os"

	"github.com/hibooboo2/tile38NukeGame/game"
)

var count int
var n = 60
var player *game.Character

func main() {
	name, _ := os.Hostname()
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	player = game.NewCharacter(name)

	ui()
	game.ClearNotifications()
}
