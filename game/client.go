package game

import (
	"log"
	"sync"

	"github.com/hibooboo2/tile38NukeGame/game/model"
)

var (
	chars    = map[string]*Character{}
	charLock = sync.Mutex{}
	events   = make(chan model.Thing)
)

// const Tile38ServerURL = "http://10.14.12.11:9851"
const Tile38ServerURL = "http://localhost:9851"

func init() {
	go func() {
		for thing := range events {
			charLock.Lock()
			c, ok := chars[thing.Nearby.ID]
			charLock.Unlock()
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
	}()
}
