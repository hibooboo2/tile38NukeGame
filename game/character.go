package game

import (
	"fmt"
	"hash/fnv"
	"image/color"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/hibooboo2/tile38NukeGame/game/model"
	"github.com/hibooboo2/tile38NukeGame/ui"
)

type Character struct {
	sync.RWMutex
	name          string
	c             *Client
	posx, posy    float64
	Things        chan model.KeyedPoint
	currentThings map[string]model.KeyedPoint
	bullets       map[int]*Bullet
	bulletChan    chan *Bullet
	colors        map[string]color.RGBA
}

type Bullet struct {
	x    float64
	y    float64
	num  int
	dir  int
	made time.Time
}

func (b *Bullet) Move() {
	switch b.dir {
	case 0:
		b.x += model.Meter
	case 1:
		b.x -= model.Meter
	case 2:
		b.y += model.Meter
	case 3:
		b.y -= model.Meter
	case 4:
		b.x += model.Meter * .707
		b.y += model.Meter * .707
	case 5:
		b.x += model.Meter * .707
		b.y -= model.Meter * .707
	case 6:
		b.x -= model.Meter * .707
		b.y += model.Meter * .707
	case 7:
		b.x -= model.Meter * .707
		b.y -= model.Meter * .707
	}
}

func NewCharacter(name string) *Character {
	c := &Character{
		Things:        make(chan model.KeyedPoint),
		currentThings: make(map[string]model.KeyedPoint),
		bullets:       make(map[int]*Bullet),
		bulletChan:    make(chan *Bullet),
		name:          name,
		c:             NewClient(Tile38ServerURL, name),
	}
	go c.handleThings()
	go c.handleBullets()
	c.MoveRel(1, 1)
	charLock.Lock()
	chars[c.name] = c
	charLock.Unlock()
	return c
}

func (c *Character) Name() string {
	return c.name
}

func (c *Character) MoveRel(x, y float64) {
	c.Lock()
	c.posx += (x * model.Meter)
	c.posy += (y * model.Meter)
	c.Unlock()
}

func (c *Character) handleBullets() {
	t := time.NewTicker(time.Millisecond * 40)
	for {
		select {
		case b := <-c.bulletChan:
			c.bullets[b.num] = b
		case <-t.C:
			for _, b := range c.bullets {
				b.Move()
				c.c.post(fmt.Sprintf("SET fleet %d point %f %f", b.num, b.x, b.y))
				if time.Since(b.made) > time.Second*5 {
					c.c.post(fmt.Sprintf("DEL fleet %d", b.num))
					delete(c.bullets, b.num)
				}
			}
		}
	}
}
func (c *Character) handleThings() {
	err := c.c.Notifications(c.name)
	if err != nil {
		panic(err)
	}
	ticker := time.NewTicker(time.Millisecond * 10)
	lastx, lasty := c.posx, c.posy
	for {
		select {
		case t := <-c.Things:
			if t.Object.Type == "delete" {
				delete(c.currentThings, t.ID)
				continue
			}

			c.Lock()
			prev := c.currentThings[t.ID]
			c.currentThings[t.ID] = t
			c.Unlock()
			if t.ID == c.name {
				log.Println("WTF!", c.name, t.ID)
			}
			if prev.Object.Coordinates.String() == t.Object.Coordinates.String() {
				log.Println("Same loc!", c.name, t.ID)
			} else {
				_, err := strconv.Atoi(t.ID)
				if err != nil {
					log.Println("Someone moved:", c.name, t.ID, t.Object.Coordinates.String())
				}
			}
		case <-ticker.C:
			c.RLock()
			if lastx != c.posx || lasty != c.posy {
				err := c.c.post(fmt.Sprintf("SET fleet %s point %f %f", c.name, c.posx, c.posy))
				if err != nil {
					log.Println("failed to set position for player!", err)
				}
				lastx, lasty = c.posx, c.posy
			}
			c.RUnlock()
		}
	}
}

func (c *Character) Shoot() {
	c.bulletChan <- &Bullet{
		x:    c.posx,
		y:    c.posy,
		num:  rand.Int(),
		dir:  rand.Intn(8),
		made: time.Now(),
	}
}

func (c *Character) Render(r ui.Renderer) {
	c.RLock()
	size := 300
	// fmt.Print(c.name, " ", len(c.currentThings), " ")
	r.DrawFilledSquare(size/2, size/2, size, color.RGBA{255, 255, 0, 0})
	r.DrawFilledSquare(size/2, size/2, 4, getUserColor(c.name))

	for _, t := range c.currentThings {
		x := c.posx
		y := c.posy
		x -= t.Object.Coordinates[1]
		y -= t.Object.Coordinates[0]
		x *= (1 / model.Meter)
		y *= (1 / model.Meter)
		r.DrawFilledSquare(int(x)+size/2, int(y)+size/2, 4, getUserColor(t.ID))
	}
	c.RUnlock()
	// fmt.Println()
}

var (
	colors    = map[string]color.RGBA{}
	colorLock = sync.Mutex{}
)

func getUserColor(name string) color.RGBA {
	colorVals := hash(name)
	colorLock.Lock()
	userColor, ok := colors[name]
	if !ok {
		userColor = color.RGBA{uint8(colorVals), uint8(colorVals >> 8), uint8(colorVals >> 16), uint8(colorVals >> 24)}
		colors[name] = userColor
	}
	colorLock.Unlock()
	return userColor
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
