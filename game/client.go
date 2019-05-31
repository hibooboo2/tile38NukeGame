package game

import (
	"log"
	"sync"
)

var (
	chars    = map[string]*Character{}
	charLock = sync.Mutex{}
	events   = make(chan Thing)
)

func init() {
	go func() {
		for thing := range events {
			charLock.Lock()
			c, ok := chars[thing.Nearby.ID]
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
			charLock.Unlock()
		}
	}()
}
