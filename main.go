package main

import (
	"os"

	"github.com/hibooboo2/tile38NukeGame/game"

	"github.com/tidwall/pinhole"
)

var count int
var opts = pinhole.DefaultImageOptions
var n = 60
var player *game.Character

func main() {
	c := make(chan struct{})
	name, _ := os.Hostname()
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	player = game.NewCharacter(name)

	ui()
	<-c
}
