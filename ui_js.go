// +build js,wasm

package main

import (
	"fmt"
	"image/color"
	"syscall/js"
	"time"
)

func ui() {
	window := js.Global()

	cv := window.Get("document").Call("getElementById", "canvas").Call("getContext", "2d")
	drawFunc := getDrawBoxFunc(cv)
	for {
		time.Sleep(time.Second / 60)
		cv.Call("clearRect", 0, 0, 300, 300)

		player.MiniMap(drawFunc, 0, 0, 300)
	}
}

func getDrawBoxFunc(cv js.Value) func(x, y, size int, c color.RGBA) {
	return func(x, y, size int, c color.RGBA) {
		cv.Set("fillStyle", fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B))
		cv.Call("fillRect", int32(x-size/2), int32(y-size/2), int32(size), int32(size))
	}
}
