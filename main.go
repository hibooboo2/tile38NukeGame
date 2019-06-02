//build !js,!wasm

package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/hibooboo2/tile38NukeGame/game"
	"github.com/hibooboo2/tile38NukeGame/ui"
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
	go func() {
		events := ui.GetKeyboardEvents()
		for evt := range events {
			log.Println("Got event:", evt)
			switch evt.Key {
			case ui.K_w:
				player.MoveRel(0, 1)
			case ui.K_a:
				player.MoveRel(1, 0)
			case ui.K_s:
				player.MoveRel(0, -1)
			case ui.K_d:
				player.MoveRel(-1, 0)
			case ui.K_q:
				player.MoveRel(.707, .707)
			case ui.K_e:
				player.MoveRel(-.707, .707)
			case ui.K_z:
				player.MoveRel(.707, -.707)
			case ui.K_c:
				player.MoveRel(-.707, -.707)
			case ui.K_SPACE:
				player.Shoot()
			}
		}
	}()
	ui.Start(player)
	game.ClearNotifications()
}
