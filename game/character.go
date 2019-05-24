package game

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"math/rand"
)

const meter = .00001 / 1.111

type Character struct {
	name          string
	c             *Client
	posx, posy    float64
	Things        chan Thing
	currentThings map[string]Thing
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
		Things:        make(chan Thing),
		currentThings: make(map[string]Thing),
		move:          make(chan struct{}),
		mini:          make(chan struct{}),
		name:          name,
		c:             NewClient("http://10.14.12.11:9851"),
	}
	go c.handleThings()
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
	c.c.Notifications(c.name, "http://10.14.12.68:8081")
	for {
		select {
		case t := <-c.Things:
			if t.Nearby.ID != "" && t.Nearby.Meters < 100 {
				fmt.Println("got thing", t.ID, t.KeyedPoint.Object.Coordinates, t.Nearby.ID, t.Nearby.Object)
				c.currentThings[t.Nearby.ID] = t
			} else {
				delete(c.currentThings, t.Nearby.ID)
			}
		case c.mini <- struct{}{}:
			<-c.mini
		case c.move <- struct{}{}:
			<-c.move
			fmt.Println("moved", c.name, c.posx*(1/meter), c.posy*(1/meter))
			c.c.post(fmt.Sprintf("SET fleet %s point %f %f", c.name, c.posx, c.posy))
		}
	}
}

func (c *Character) MiniMap(img *image.RGBA, startX, startY, size int) {
	<-c.mini
	// fmt.Print(c.name, " ", len(c.currentThings), " ")
	for x := startX; x < startX+size; x++ {
		for y := startY; y < startY+size; y++ {
			img.Set(x, y, color.RGBA{255, 255, 0, 0})
		}
	}
	drawBox(img, size/2+startX, size/2+startY, 4, color.RGBA{120, 255, 100, 0})

	for _, t := range c.currentThings {
		x := c.posx
		y := c.posy
		x -= t.Nearby.Object.Coordinates[0]
		y -= t.Nearby.Object.Coordinates[1]
		x *= (1 / meter)
		y *= (1 / meter)

		drawBox(img, int(x)+size/2+startX, int(y)+size/2+startY, 4, c.GetUserColor(t.Nearby.ID))
	}
	// fmt.Println()
	c.mini <- struct{}{}
}
func (c *Character) GetUserColor(name string) color.RGBA {
	if c.colors == nil {
		c.colors = make(map[string]color.RGBA)
	}
	userColor, ok := c.colors[name]
	if !ok {
		userColor = color.RGBA{255, uint8(rand.Int() % 255), uint8(rand.Int() % 255), uint8(rand.Int() % 255)}
		c.colors[name] = userColor
	}
	return userColor
}

func drawBox(img *image.RGBA, x, y, size int, c color.RGBA) {
	b := img.Bounds()
	if x > b.Min.X && x < b.Max.X && y > b.Min.Y && y < b.Max.Y {
		for i := -size / 2; i < size/2; i++ {
			for j := -size / 2; j < size/2; j++ {
				img.Set(x+i, y+j, c)
			}
		}
	}
}