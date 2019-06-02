package game

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"image/color"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/hibooboo2/tile38NukeGame/ui"
)

const meter = .00001 / 1.111

type Character struct {
	sync.RWMutex
	name          string
	c             *Client
	posx, posy    float64
	Things        chan KeyedPoint
	currentThings map[string]KeyedPoint
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
		b.x += meter
	case 1:
		b.x -= meter
	case 2:
		b.y += meter
	case 3:
		b.y -= meter
	case 4:
		b.x += meter * .707
		b.y += meter * .707
	case 5:
		b.x += meter * .707
		b.y -= meter * .707
	case 6:
		b.x -= meter * .707
		b.y += meter * .707
	case 7:
		b.x -= meter * .707
		b.y -= meter * .707
	}
}

type Thing struct {
	KeyedPoint
	Command string `json:"command"`
	Group   string `json:"group"`
	Detect  string `json:"detect"`
	Hook    string `json:"hook"`
	Time    string `json:"time"`
	// Faraway DistancePoint `json:"faraway"`
	Nearby DistancePoint `json:"nearby"`
}

type KeyedPoint struct {
	Key    string `json:"key"`
	ID     string `json:"id"`
	Object Point  `json:"object"`
}

type DistancePoint struct {
	KeyedPoint
	Meters float64 `json:"meters"`
}
type Point struct {
	Type        string `json:"type"`
	Coordinates Coord  `json:"coordinates"`
}

type Coord []float64

func (c Coord) String() string {
	var b bytes.Buffer
	b.WriteString("[ ")
	for i, val := range c {
		b.WriteString(fmt.Sprintf("%f", val*(1/meter)))
		if i < len(c)-1 {
			b.WriteString(" ,")
		}
	}
	b.WriteString(" ]")
	return b.String()
}

func NewCharacter(name string) *Character {
	c := &Character{
		Things:        make(chan KeyedPoint),
		currentThings: make(map[string]KeyedPoint),
		bullets:       make(map[int]*Bullet),
		bulletChan:    make(chan *Bullet),
		name:          name,
		c:             NewClient(Tile38ServerURL),
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
	c.posx += (x * meter)
	c.posy += (y * meter)
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
		x *= (1 / meter)
		y *= (1 / meter)
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
