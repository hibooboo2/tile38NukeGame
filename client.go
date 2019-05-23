package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/tidwall/pinhole"
)

var (
	chars    = map[string]*Character{}
	charLock = sync.Mutex{}
)

func players() {
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		thing := Thing{}
		err := json.NewDecoder(req.Body).Decode(&thing)
		if err != nil {
			fmt.Println(err)
		}

		charLock.Lock()
		c, ok := chars[thing.ID]
		charLock.Unlock()
		if ok {
			c.things <- thing
		} else {
			log.Printf("%s %#v", req.URL.String(), thing)
		}
		req.Body.Close()
		resp.WriteHeader(http.StatusOK)
	})
	go http.ListenAndServe(":8081", nil)
	c := &Client{&http.Client{Timeout: time.Second * 2}, "none"}
	c.post("PDELHOOK *")
	go MoveAround("james")
	go MoveAround("george")
	go MoveAround("jeff")
}

const meter = .00001 / 1.111

type Client struct {
	c    *http.Client
	name string
}
type Character struct {
	things        chan Thing
	currentThings map[string]Thing
	wait          chan struct{}
}

func (c *Character) HandleThings() {
	for {
		select {
		case t := <-c.things:
			if t.Nearby.ID != "" {
				c.currentThings[t.Nearby.ID] = t
			}
		case c.wait <- struct{}{}:
		}
	}
}

func (c *Character) GetPinHole() *pinhole.Pinhole {
	<-c.wait
	p := pinhole.New()
	for _, t := range c.currentThings {

		x := t.Nearby.Object.Coordinates[0] * (1 / meter)
		y := t.Nearby.Object.Coordinates[1] * (1 / meter)
		log.Println(t.KeyedPoint.Object.Coordinates, t.Nearby.Object.Coordinates, x, y)

		p.DrawDot(x, y, 0, .1)
	}
	return p
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
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

func MoveAround(name string) {
	person := Character{make(chan Thing), make(map[string]Thing), make(chan struct{})}
	charLock.Lock()
	chars[name] = &person
	charLock.Unlock()
	go person.HandleThings()

	c := &Client{&http.Client{Timeout: time.Second * 2}, name}
	c.Notifications()
	for {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
		c.post(fmt.Sprintf("set fleet %s point %f %f", c.name, float64(rand.Intn(10))*meter, float64(rand.Intn(10))*meter))
	}
}

func (c *Client) Notifications() {
	c.post(fmt.Sprintf("SETHOOK %[1]s http://127.0.0.1:8081/%[1]s NEARBY fleet MATCH %[1]s FENCE ROAM fleet * 1000", c.name))
}

func (c *Client) post(cmd string) {
	resp, err := c.c.Post("http://localhost:9851", "", strings.NewReader(cmd))
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	io.Copy(os.Stderr, resp.Body)
}

func (c *Client) get(cmd string) {
	resp, err := c.c.Post("http://localhost:9851", "", strings.NewReader(cmd))
	if err != nil {
		fmt.Println(err)
		return
	}
	io.Copy(os.Stderr, resp.Body)
}
