package main

import (
	"blog/drawer"
	"image"
	"image/color"
	"log"
)

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

	// draw a black frame
	p.DrawRectangle(image.Rect(0, 0, 300, 300), 25, color.Black)

	// draw a white circle in the center
	p.DrawFilledCircle(drawer.Cc(150, 150, 25), color.White)

	// draw a little red dot in the center
	p.DrawFilledCircle(drawer.Cc(150, 150, 10), color.RGBA{R: 255, G: 0, B: 0, A: 255})

	if err := p.SavePNG("amazing_logo.png"); err != nil {
		log.Fatalln(err)
	}
}
