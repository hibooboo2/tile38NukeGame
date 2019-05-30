package main

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"runtime"
)

func init() {
	if runtime.GOOS != "windows" {
		// log.Println("Not windows dling anyways")
		return
	}
	inf, err := os.Lstat("SDL2.dll")
	if os.IsExist(err) {
		return
	}

	resp, err := http.Get("http://libsdl.org/release/SDL2-2.0.9-win32-x64.zip")
	if err != nil {
		panic(err)
	}
	f, err := os.Create("sdl.zip")
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		panic(err)
	}
	f.Close()
	resp.Body.Close()

	f, err = os.Open("sdl.zip")
	if err != nil {
		panic(err)
	}
	inf, _ = f.Stat()
	r, err := zip.NewReader(f, inf.Size())
	if err != nil {
		panic(err)
	}
	for _, f := range r.File {
		if f.Name == "SDL2.dll" {
			r, err := f.Open()
			if err != nil {
				panic(err)
			}
			f, err := os.Create("SDL2.dll")
			if err != nil {
				panic(err)
			}
			_, err = io.Copy(f, r)
			if err != nil {
				panic(err)
			}
			f.Close()
			r.Close()
		}
	}
	f.Close()
	os.Remove("sdl.zip")
}
