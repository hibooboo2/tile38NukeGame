// +build js,wasm

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"syscall/js"
	"time"
)

func events() {
	//Code to be able to connect to a websocket and connect to a server
	ws := js.Global().Get("WebSocket").New("ws://localhost:8000/ws")
	events := make(chan string)
	ws.Call("addEventListener", "open", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("open")
		ws.Set("onmessage", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			log.Println(this.String())
			msg := args[0]
			data := msg.Get("data")
			events <- data.String()
			return nil
		}))

		someObj := map[string]interface{}{
			"boom":   3,
			"hello!": 324.234,
			"sweet":  false,
		}
		data, _ := json.Marshal(someObj)
		log.Println(ws.Call("send", string(data)))
		return nil
	}))

	go func() {
		last := time.Now()
		for evt := range events {
			log.Println(time.Since(last).Nanoseconds() / int64(time.Millisecond))
			last = time.Now()
			vals := map[string]interface{}{}
			json.Unmarshal([]byte(evt), &vals)
			log.Println(vals)
		}
	}()
}
