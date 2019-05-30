// +build !nuke

package main

import (
	"fmt"
	"image/color"
	"log"
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func ui() {
	r, cleanup := getRenderer(300, 300)
	defer cleanup()

	EventLoop(r)
}

func getRenderer(h, w int32) (*sdl.Renderer, func()) {
	runtime.LockOSThread()

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	if err := ttf.Init(); err != nil {
		panic(err)
	}

	window, r, err := sdl.CreateWindowAndRenderer(w, h, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	window.SetResizable(true)
	window.SetBordered(true)
	// window.SetGrab(true)
	// window.SetWindowOpacity(0.4)

	return r,
		func() {
			window.Destroy()
			sdl.Quit()
			ttf.Quit()
			runtime.UnlockOSThread()
		}

}

var (
	events = make(chan sdl.Event, 2)
	quit   = make(chan struct{})
)

func EventLoop(r *sdl.Renderer) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	go PaintLoop(r)

	go func() {
		for event := range events {
			switch e := event.(type) {
			case *sdl.KeyboardEvent:
				log.Println(e.Keysym)
				switch e.Keysym.Scancode {
				case sdl.SCANCODE_W:
					player.MoveRel(0, 1)
				case sdl.SCANCODE_S:
					player.MoveRel(0, -1)
				case sdl.SCANCODE_A:
					player.MoveRel(1, 0)
				case sdl.SCANCODE_D:
					player.MoveRel(-1, 0)
				}
			case *sdl.MouseMotionEvent:
			case *sdl.MouseButtonEvent:
			case *sdl.QuitEvent:
				close(quit)
				return
			default:
				log.Printf("%T", e)
			}
		}
	}()

	for {
		select {
		case events <- sdl.WaitEvent():
			// time.Sleep(time.Millisecond)
		case <-quit:
			close(events)
			fmt.Println("Quitting")
			return
		}
	}
}

func PaintLoop(r *sdl.Renderer) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ticker := time.NewTicker(time.Second / 60)
	for {
		select {
		case <-quit:
			return
		case <-ticker.C:
		}
		// r.Clear()
		player.MiniMap(getDrawBoxFunc(r), 0, 0, 300)
		r.Present()
	}
}

func getDrawBoxFunc(r *sdl.Renderer) func(x, y, size int, c color.RGBA) {
	return func(x, y, size int, c color.RGBA) {
		r.SetDrawColor(c.R, c.G, c.B, c.A)
		r.FillRect(&sdl.Rect{int32(x - size/2), int32(y - size/2), int32(size), int32(size)})
	}
}
