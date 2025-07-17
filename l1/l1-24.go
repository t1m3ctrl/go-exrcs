package main

import (
	"fmt"
	"math"
)

type Point struct {
	x, y float64
}

func NewPoint(x, y float64) (p *Point) {
	return &Point{
		x: x,
		y: y,
	}
}

func (p *Point) Distance(other *Point) float64 {
	return math.Sqrt((p.x-other.x)*(p.x-other.x) + (p.y-other.y)*(p.y-other.y))
}

func main() {
	point := NewPoint(0, 0)
	other := NewPoint(0, 5)
	fmt.Println(point.Distance(other))
}
