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
)

const meter = .00001 / 1.111

type Character struct {
	name          string
	c             *Client
	posx, posy    float64
	Things        chan KeyedPoint
	currentThings map[string]KeyedPoint
	bullets       map[int]*Bullet
	bulletChan    chan *Bullet
	move          chan struct{}
	mini          chan struct{}
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
		move:          make(chan struct{}),
		mini:          make(chan struct{}),
		name:          name,
		c:             NewClient("http://10.14.12.11:9851"),
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
	<-c.move
	c.posx += (x * meter)
	c.posy += (y * meter)
	c.move <- struct{}{}
}

func (c *Character) handleBullets() {
	t := time.NewTicker(time.Millisecond * 40)
	for {
		select {
		case b := <-c.bulletChan:
			c.bullets[b.num] = b
		case <-t.C:
			for _, b := range c.bullets {
				c.c.post(fmt.Sprintf("SET fleet %d point %f %f", b.num, b.x, b.y))
				b.Move()
				if time.Since(b.made) > time.Second*5 {
					c.c.post(fmt.Sprintf("DEL fleet %d", b.num))
					delete(c.bullets, b.num)
				}
			}
		}
	}
}
func (c *Character) handleThings() {
	c.c.Notifications(c.name)
	ticker := time.NewTicker(time.Millisecond * 1)
	lastx, lasty := c.posx, c.posy
	for {
		select {
		case t := <-c.Things:
			if t.Object.Type == "delete" {
				delete(c.currentThings, t.ID)
				continue
			}

			prev := c.currentThings[t.ID]
			c.currentThings[t.ID] = t

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
		case c.mini <- struct{}{}:
			<-c.mini
		case c.move <- struct{}{}:
			<-c.move
			fmt.Println("moved", c.name, c.posx*(1/meter), c.posy*(1/meter))
		case <-ticker.C:
			if lastx != c.posx || lasty != c.posy {
				ok := c.c.post(fmt.Sprintf("SET fleet %s point %f %f", c.name, c.posx, c.posy))
				if !ok {
					log.Println(ok, c.posx, c.posy)
				}
				lastx, lasty = c.posx, c.posy
			}
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
