package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/hibooboo2/tile38NukeGame/game"
)

var (
	chars    = map[string]*game.Character{}
	charLock = sync.Mutex{}
)

func playersWebHooks() {
	log.SetFlags(log.Lshortfile)
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		thing := game.Thing{}
		err := json.NewDecoder(req.Body).Decode(&thing)
		if err != nil {
			fmt.Println(err)
		}

		charLock.Lock()
		c, ok := chars[thing.ID]
		charLock.Unlock()
		if ok {
			c.Things <- thing
		} else {
			log.Printf("%s %#v", req.URL.String(), thing.ID, thing.Nearby.ID)
		}
		req.Body.Close()
		resp.WriteHeader(http.StatusOK)
	})
	go http.ListenAndServe(":8081", nil)
	time.Sleep(time.Second)
}
