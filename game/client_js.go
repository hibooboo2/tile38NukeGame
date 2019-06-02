// +build js,wasm

package game

import (
	"encoding/json"
	"log"
	"strings"
	"syscall/js"

	"github.com/pkg/errors"
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

func (c *Client) post(cmd string) error {
	cmd = "tile38: " + cmd
	go c.ws.Call("send", cmd)
	resp := <-c.data
	log.Println(resp)
	if strings.Contains(resp, "true") {
		return nil
	}
	return errors.Errorf("failed to send command")
}

func (c *Client) Notifications(name string) error {
	return nil
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
