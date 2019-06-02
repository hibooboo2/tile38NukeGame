package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(*http.Request) bool { return true },
} // use default options

func main() {
	//Basic server that serves up js ui and also exposes a endpoint for websocket to connect to.
	http.HandleFunc("/events", echo)

	http.Handle("/", http.FileServer(http.Dir(".")))
	http.ListenAndServe(":8000", nil)
}

// Hold players that are currently connected
// Send players updates when players move
// Store player positions...
// Notifications for players!

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	// Read all messages and log them
	// Write all messages that we get from the webhook we make
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		msg := string(message)
		if strings.HasPrefix(msg, "tile38: ") {
			resp, err := http.Post("http://localhost:9851/", "", strings.NewReader(strings.TrimPrefix(msg, "tile38: ")))
			if err != nil {
				c.WriteMessage(mt, []byte(err.Error()))
			} else {
				data, _ := ioutil.ReadAll(resp.Body)
				err = c.WriteMessage(mt, data)
				if err != nil {
					log.Println("write:", err)
					break
				}
				log.Println(string(data))
				resp.Body.Close()
			}
		} else {
			log.Printf("recv: %s", message)
			err = c.WriteMessage(mt, message)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}
}
