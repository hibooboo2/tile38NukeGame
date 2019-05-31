package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

func main() {
	//Basic server that serves up js ui and also exposes a endpoint for websocket to connect to.
	http.HandleFunc("/ws", echo)

	http.Handle("/", http.FileServer(http.Dir(".")))
	http.ListenAndServe(":8000", nil)
}

// Hold players that are currently connected
// Send players updates when players move
// Store player positions...

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	go func() {
		i := 0
		for {
			i++
			time.Sleep(time.Millisecond * 100)
			c.WriteJSON(map[string]string{
				"hello!": "boy!",
				"count":  fmt.Sprintf("%d", i),
			})
		}
	}()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
