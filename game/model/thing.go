package model

import (
	"bytes"
	"fmt"
)

const Meter = .00001 / 1.111

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
		b.WriteString(fmt.Sprintf("%f", val*(1/Meter)))
		if i < len(c)-1 {
			b.WriteString(" ,")
		}
	}
	b.WriteString(" ]")
	return b.String()
}
