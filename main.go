package main

import (
	"fmt"
	"math"
	"math/rand"
	"syscall/js"
)

// Gravity is the acceleration of the object to the bottom of the screen.
const gravity int = 1

func main() {
	done := make(chan struct{}, 0)
	doc := js.Global().Get("document")
	body := doc.Get("body")
	width := body.Get("clientWidth").Int()
	height := body.Get("clientHeight").Int()

	canvas := doc.Call("getElementById", "canvas")

	ball := &Ball{
		Radius:   height / 10,
		Position: NewPoint(width/2, height/2),
		Velocity: NewVelocity(rand.Intn(10)-5, rand.Intn(10)-5),
		Color:    RandomColor(),
	}

	canvas.Set("width", width)
	canvas.Set("height", height)

	context := canvas.Call("getContext", "2d")

	var renderFrame js.Func
	renderFrame = js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		context.Call("clearRect", 0, 0, width, height)

		ball.Update()

		if ball.Position.x-ball.Radius < 0 {
			ball.Bounce(Left)
		}
		if ball.Position.x+ball.Radius > width {
			ball.Bounce(Right)
		}
		if ball.Position.y-ball.Radius < 0 {
			ball.Bounce(Top)
		}
		if ball.Position.y+ball.Radius > height {
			ball.Bounce(Bottom)
		}

		context.Set("fillStyle", ball.Color.Hex())
		context.Call("beginPath")
		context.Call("arc", ball.Position.x, ball.Position.y, ball.Radius, 0, 2*math.Pi)
		context.Call("fill")

		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})
	defer renderFrame.Release()

	fmt.Println("Rendering...")

	js.Global().Call("requestAnimationFrame", renderFrame)
	<-done
}

type Point struct {
	x int
	y int
}

func NewPoint(x, y int) *Point {
	return &Point{x, y}
}

type Velocity struct {
	x int
	y int
}

func NewVelocity(x, y int) *Velocity {
	return &Velocity{x, y}
}

type Color struct {
	// Rgb is the color of the object. The largest byte is unused.
	Rgb uint32
}

func RandomColor() *Color {
	rgb := rand.Intn(0xFFFFFF)
	return &Color{Rgb: uint32(rgb)}
}

func (c Color) Hex() string {
	return fmt.Sprintf("#%06X", c.Rgb)
}

func (c Color) RGBA() (r, g, b, a uint32) {
	r = (c.Rgb >> 16) & 0xFF
	g = (c.Rgb >> 8) & 0xFF
	b = c.Rgb & 0xFF
	a = 0xFF
	return
}

// Shift "increments" the hue by 1, wrapping around at 0xFFFFFF.
func (c *Color) Shift(amount uint32) {
	c.Rgb = (c.Rgb + amount) % 0xFFFFFF
}

type Ball struct {
	// Radius is the size of the object.
	Radius int
	// Position is the current position of the object.
	Position *Point
	// Velocity is the speed of the object.
	Velocity *Velocity
	// Color is the color of the object.
	Color *Color
}

type Boundary = uint8

const (
	Top Boundary = iota
	Right
	Bottom
	Left
)

func (b *Ball) Update() {
	b.Color.Shift(0b100)
	b.Position.x += b.Velocity.x
	b.Position.y += b.Velocity.y
	b.Velocity.y += gravity
}

func (b *Ball) Bounce(boundary Boundary) {
	switch boundary {
	case Top:
		b.Velocity.y = -b.Velocity.y
		b.Position.y += b.Radius
	case Right:
		b.Velocity.x = -b.Velocity.x
		b.Position.x -= b.Radius
	case Bottom:
		b.Velocity.y = -b.Velocity.y
		b.Position.y -= b.Radius
	case Left:
		b.Velocity.x = -b.Velocity.x
		b.Position.x += b.Radius
	}
}
