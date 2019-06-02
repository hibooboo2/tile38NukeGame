package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/hibooboo2/tile38NukeGame/game/model"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(*http.Request) bool { return true },
} // use default options

func main() {
	log.SetFlags(log.Lshortfile)
	//Basic server that serves up js ui and also exposes a endpoint for websocket to connect to.
	http.HandleFunc("/events", echo)
	http.HandleFunc("/tile38", tile38)
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.ListenAndServe(":8000", nil)
}

var users = map[string]*websocket.Conn{}
var usersLock sync.Mutex

func tile38(resp http.ResponseWriter, req *http.Request) {
	log.Println("Tile 38 hooked")
	thing := model.Thing{}
	err := json.NewDecoder(req.Body).Decode(&thing)
	if err != nil {
		log.Println(err)
		return
	}
	usersLock.Lock()
	con, ok := users[thing.Nearby.ID]
	if ok {
		con.WriteJSON(thing)
	} else {
		if thing.Nearby.ID == "" {
			thing.KeyedPoint.Object.Type = "delete"
			for _, con := range users {
				con.WriteJSON(thing)
			}
		} else {
			log.Printf("ID [%s] nearID [%s]", thing.ID, thing.Nearby.ID)
		}
	}
	usersLock.Unlock()

	req.Body.Close()
	resp.WriteHeader(http.StatusOK)
}

// Hold players that are currently connected
// Send players updates when players move
// Store player positions...
// Notifications for players!

func echo(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	usersLock.Lock()
	users[id] = c
	usersLock.Unlock()
	resp, err := http.Post("http://localhost:9851/", "", strings.NewReader(fmt.Sprintf("SETHOOK %[1]s %[2]s%[1]s NEARBY fleet FENCE ROAM fleet %[1]s 1000", id, "http://localhost:8000/tile38?id=")))
	if err != nil {
		log.Println("Failed to setup webhook!")
	} else {
		io.Copy(os.Stdout, resp.Body)
		resp.Body.Close()
		log.Println("Setup webhook for: ", id)
	}
	defer func() {
		usersLock.Lock()
		delete(users, id)
		usersLock.Unlock()
	}()
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
				// log.Println(string(data))
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
