//build !js,!wasm

package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/hibooboo2/tile38NukeGame/game"
)

var count int
var n = 60
var player *game.Character

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.Lshortfile)
	name, _ := os.Hostname()
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	player = game.NewCharacter(name)
	log.Println("Starting game!")
	ui()
	game.ClearNotifications()
}
