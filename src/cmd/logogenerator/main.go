package main

import (
	"blog/drawer"
	"image"
	"image/color"
	"log"
	"math/rand"
)

//go:generate go run main.go

// source: https://colorswall.com/palette/102
var rainbowColors = []color.Color{
	color.RGBA{R: 255, G: 0, B: 0, A: 255},     // red
	color.RGBA{R: 255, G: 165, B: 0, A: 255},   // orange
	color.RGBA{R: 255, G: 255, B: 0, A: 255},   // yellow
	color.RGBA{R: 0, G: 128, B: 0, A: 255},     // green
	color.RGBA{R: 0, G: 0, B: 255, A: 255},     // blue
	color.RGBA{R: 75, G: 0, B: 130, A: 255},    // indigo
	color.RGBA{R: 238, G: 130, B: 238, A: 255}, // violet
}

func main() {
	p := drawer.NewPicture(image.Rect(0, 0, 300, 300))

	// draw insanity
	for i := 300; i > 0; i-- {
		var c color.Color
		if i%2 == 0 {
			c = color.Black
		} else {
			c = color.White
		}
		p.DrawFilledCircle(drawer.Cc(150, 150, float64(i)), c)
	}

	// draw 1000 random filled circles with random rainbow colors
	for i := 0; i < 1000; i++ {
		x := 0 + rand.Intn(300)
		y := 0 + rand.Intn(300)
		r := 1 + rand.Float64()*5
		c := rainbowColors[rand.Intn(len(rainbowColors))]
		p.DrawFilledCircle(drawer.Cc(x, y, r), c)
	}

	// draw a rainbow circle in the center
	currentColor := 0
	for i := 50; i > 0; i-- {
		p.DrawFilledCircle(drawer.Cc(150, 150, float64(i)), rainbowColors[currentColor%len(rainbowColors)])
		currentColor++
	}

	// draw transparency around the picture
	p.DrawCircle(drawer.Cc(150, 150, 200), 100, color.Transparent)

	// draw black round frame around
	p.DrawCircle(drawer.Cc(150, 150, 145), 10, color.Black)

	if err := p.SavePNG("amazing_logo.png"); err != nil {
		log.Fatalln(err)
	}
}
