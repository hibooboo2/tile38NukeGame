// +build js,wasm

package ui

import (
	"fmt"
	"image/color"
	"log"
	"syscall/js"
	"time"
)

var canvasContext js.Value

func Start(r Renderable) {
	go events()
	window := js.Global()
	canvasContext = window.Get("document").Call("getElementById", "canvas").Call("getContext", "2d")
	ren := &JsRenderer{canvasContext}
	ticker := time.NewTicker(time.Second / 60)
	for {
		select {
		// case <-quit:
		// 	return
		case <-ticker.C:
		}
		canvasContext.Call("clearRect", 0, 0, 300, 300)
		r.Render(ren)
	}
}

func events() {
	keydown := make(chan rune, 2)
	keyup := make(chan rune, 2)
	js.Global().Get("document").Set("onkeydown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		log.Printf("evt: %T %s", args[0], args[0].Get("key"))
		if args[0].Get("repeat").String() == "false" {
			key := args[0].Get("key").String()[0]
			keydown <- rune(key)
		}
		return nil
	}))
	js.Global().Get("document").Set("onkeyup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		log.Printf("evt: %T %s", args[0], args[0].Get("key"))
		key := args[0].Get("key").String()[0]
		keyup <- rune(key)
		return nil
	}))
	go func() {
		keymap := map[rune]bool{}
		t := time.NewTicker(time.Millisecond * 10)
		for {
			select {
			case key := <-keydown:
				keymap[key] = true
			case key := <-keyup:
				delete(keymap, key)
			case <-t.C:
				for key := range keymap {
					mainKeyBoardEvents <- KeyboardEvent{
						Key: keyMap[key],
					}
				}
			}
		}
	}()
}

type JsRenderer struct {
	cv js.Value
}

var _ Renderer = &JsRenderer{}

func (r *JsRenderer) DrawFilledSquare(x, y, size int, c color.RGBA) {
	r.cv.Set("fillStyle", fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B))
	r.cv.Call("fillRect", int32(x-size/2), int32(y-size/2), int32(size), int32(size))
}
