// +build !js,!wasm

package ui

import (
	"fmt"
	"image/color"
	"log"
	"runtime"
	"time"

	sdl "github.com/veandco/go-sdl2/sdl"
)

func init() {

}

func ui(ren Renderable) {
	r, cleanup := getRenderer(300, 300)
	defer cleanup()

	EventLoop(r, ren)
}

func getRenderer(h, w int32) (*sdl.Renderer, func()) {
	runtime.LockOSThread()

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	// if err := ttf.Init(); err != nil {
	// 	panic(err)
	// }

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
			// ttf.Quit()
			runtime.UnlockOSThread()
		}

}

func EventLoop(r *sdl.Renderer, ren Renderable) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	log.Println(keyMap[0])
	events := make(chan sdl.Event, 2)
	quit := make(chan struct{})
	sdl.GetKeyboardState()
	go func() {
		for event := range events {
			switch e := event.(type) {
			case *sdl.KeyboardEvent:
				k := keyMap[int32(e.Keysym.Sym)]
				mainKeyBoardEvents <- KeyboardEvent{
					Key: k,
				}
			case *sdl.TextInputEvent:
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
	go PaintLoop(r, ren)

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

func Start(r Renderable) {
	// go func() {
	// 	for event := range events {
	// 		switch e := event.(type) {
	// 		case *sdl.KeyboardEvent:
	// 			log.Println(e.Keysym)
	// 			switch e.Keysym.Scancode {
	// 			case sdl.SCANCODE_W:
	// 				player.MoveRel(0, 1)
	// 			case sdl.SCANCODE_S:
	// 				player.MoveRel(0, -1)
	// 			case sdl.SCANCODE_A:
	// 				player.MoveRel(1, 0)
	// 			case sdl.SCANCODE_D:
	// 				player.MoveRel(-1, 0)
	// 			case sdl.SCANCODE_SPACE:
	// 				// if e.Repeat > 0 {
	// 				// 	continue
	// 				// }
	// 				player.Shoot()
	// 			}
	// 		case *sdl.MouseMotionEvent:
	// 		case *sdl.MouseButtonEvent:
	// 		case *sdl.QuitEvent:
	// 			close(quit)
	// 			return
	// 		default:
	// 			log.Printf("%T", e)
	// 		}
	// 	}
	// }()
	ui(r)
}

func PaintLoop(r *sdl.Renderer, ren Renderable) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	u := &SdlRenderer{r}
	ticker := time.NewTicker(time.Second / 60)
	for {
		select {
		// case <-quit:
		// 	return
		case <-ticker.C:
		}
		r.Clear()
		ren.Render(u)
		// player.MiniMap(getDrawBoxFunc(r), 0, 0, 300)
		r.Present()
	}
}

type SdlRenderer struct {
	r *sdl.Renderer
}

var _ Renderer = &SdlRenderer{}

func (r *SdlRenderer) DrawFilledSquare(x, y, size int, c color.RGBA) {
	r.r.SetDrawColor(c.R, c.G, c.B, c.A)
	r.r.FillRect(&sdl.Rect{int32(x - size/2), int32(y - size/2), int32(size), int32(size)})
}
