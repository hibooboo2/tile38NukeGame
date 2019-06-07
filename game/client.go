package game

import (
	"log"
	"sync"
	"time"

	"github.com/hibooboo2/tile38NukeGame/game/model"
)

var (
	chars    = map[string]*Character{}
	charLock = sync.Mutex{}
	events   = make(chan model.Thing, 2)
)

// const Tile38ServerURL = "http://10.14.12.11:9851"
const Tile38ServerURL = "http://localhost:9851"

func init() {
	go func() {
		t := time.NewTicker(time.Millisecond * 4)
		for {
			select {
			case <-t.C:
			case thing := <-events:
				log.Println("Got thing!")
				charLock.Lock()
				c, ok := chars[thing.Nearby.ID]
				charLock.Unlock()
				log.Println("Past Lock!")
				if ok {
					c.Things <- thing.KeyedPoint
				} else {
					if thing.Nearby.ID == "" {
						thing.KeyedPoint.Object.Type = "delete"
						for _, c := range chars {
							c.Things <- thing.KeyedPoint
						}
					} else {
						log.Printf("ID [%s] nearID [%s]", thing.ID, thing.Nearby.ID)
					}
				}
			}
		}
	}()
}
