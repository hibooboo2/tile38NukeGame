package game

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"image/color"
	"log"
	"sync"
	"time"
)

const meter = .00001 / 1.111

type Character struct {
	name          string
	c             *Client
	posx, posy    float64
	Things        chan KeyedPoint
	currentThings map[string]KeyedPoint
	move          chan struct{}
	mini          chan struct{}
	colors        map[string]color.RGBA
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
		move:          make(chan struct{}),
		mini:          make(chan struct{}),
		name:          name,
		c:             NewClient("http://10.14.12.11:9851"),
	}
	go c.handleThings()
	charLock.Lock()
	chars[c.name] = c
	charLock.Unlock()
	return c
}

func (c *Character) Name() string {
	return c.name
}

func (c *Character) MoveRel(x, y float64) {
	<-c.move
	c.posx += (x * meter)
	c.posy += (y * meter)
	c.move <- struct{}{}
}

func (c *Character) handleThings() {
	c.c.Notifications(c.name)
	i := 0
	ticker := time.NewTicker(time.Millisecond * 100)
	lastx, lasty := c.posx, c.posy
	for {
		select {
		case t := <-c.Things:
			i++
			prev := c.currentThings[t.ID]
			c.currentThings[t.ID] = t
			if t.ID == c.name {
				log.Println("WTF!", c.name, i)
			}
			if prev.Object.Coordinates.String() == t.Object.Coordinates.String() {
				log.Println("Same loc!", c.name, t.ID)
			} else {
				log.Println("Someone moved:", c.name, t.ID, t.Object.Coordinates.String())
			}
		case c.mini <- struct{}{}:
			<-c.mini
		case c.move <- struct{}{}:
			<-c.move
			fmt.Println("moved", c.name, c.posx*(1/meter), c.posy*(1/meter))
		case <-ticker.C:
			if lastx != c.posx || lasty != c.posy {
				c.c.post(fmt.Sprintf("SET fleet %s point %f %f", c.name, c.posx, c.posy))
				lastx, lasty = c.posx, c.posy
			}
		}
	}
}

func (c *Character) MiniMap(drawBox func(x, y, size int, c color.RGBA), startX, startY, size int) {
	<-c.mini
	// fmt.Print(c.name, " ", len(c.currentThings), " ")
	drawBox(startX+size/2, startY+size/2, size, color.RGBA{255, 255, 0, 0})
	drawBox(size/2+startX, size/2+startY, 4, getUserColor(c.name))

	for _, t := range c.currentThings {
		x := c.posx
		y := c.posy
		x -= t.Object.Coordinates[1]
		y -= t.Object.Coordinates[0]
		x *= (1 / meter)
		y *= (1 / meter)

		drawBox(int(x)+size/2+startX, int(y)+size/2+startY, 4, getUserColor(t.ID))
	}
	// fmt.Println()
	c.mini <- struct{}{}
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
