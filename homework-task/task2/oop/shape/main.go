package main

import (
	"fmt"
	"math"
)

func main() {

	var s Shape

	s = Rectangle{Width: 20.56, Height: 54.55}
	fmt.Printf("rectangle,area:%.2f,perimeter:%.2f", s.Area(), s.Perimeter())
	fmt.Println()

	s = Circle{Radius: 20.56}
	fmt.Printf("circle,area:%.2f,perimeter:%.2f", s.Area(), s.Perimeter())
}

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

type Circle struct {
	Radius float64
}

func (r Rectangle) Area() float64 {
	fmt.Println("call Rectangle Area function")
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	fmt.Println("call Rectangle Perimeter function")
	return 2 * (r.Width + r.Height)
}

func (c Circle) Area() float64 {
	fmt.Println("call Circle Area function")
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	fmt.Println("call Circle Perimeter function")
	return 2 * math.Pi * c.Radius

}
