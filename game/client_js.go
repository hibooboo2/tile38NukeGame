// +build js,wasm

package game

import (
	"encoding/json"
	"strings"
	"syscall/js"
)

type Client struct {
	ws   js.Value
	data chan string
}

func NewClient(baseUrl string) *Client {
	c := &Client{}
	c.ws = js.Global().Get("WebSocket").New("ws://localhost:8000/events")
	started := make(chan struct{})
	c.data = make(chan string)
	c.ws.Call("addEventListener", "open", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		c.ws.Set("onmessage", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			c.data <- args[0].Get("data").String()
			return nil
		}))
		close(started)
		return nil
	}))
	<-started
	return c
}

func (c *Client) post(cmd string) bool {
	cmd = "tile38: " + cmd
	c.ws.Call("send", cmd)
	resp := <-c.data
	return strings.Contains(resp, "true")
}

func (c *Client) Notifications(name string) bool {
	return false
}

func ClearNotifications() {

}
func init() {
	go func() {
		ws := js.Global().Get("WebSocket").New("ws://localhost:8000/events")
		ws.Call("addEventListener", "open", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			ws.Set("onmessage", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				evt := Thing{}
				if json.Unmarshal([]byte(args[0].Get("data").String()), &evt) == nil {
					events <- evt
				}
				return nil
			}))
			return nil
		}))
	}()
}
