package drawer

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
)

type Image interface {
	Bounds() image.Rectangle
	At(x, y int) color.Color

	DrawPixel(x, y int, c color.Color)

	DrawCircle(cc Circle, w float64, c color.Color)
	DrawFilledCircle(cc Circle, c color.Color)

	DrawRectangle(r image.Rectangle, w float64, c color.Color)
	DrawFilledRectangle(r image.Rectangle, c color.Color)
}

type Circle struct {
	Center image.Point
	Radius float64
}

func Cc(x, y int, r float64) Circle {
	return Circle{
		Center: image.Pt(x, y),
		Radius: r,
	}
}

type Picture struct {
	img image.RGBA
}

func NewPicture(bounds image.Rectangle) *Picture {
	return &Picture{
		img: *image.NewRGBA(bounds),
	}
}

func (p *Picture) Bounds() image.Rectangle {
	return p.img.Bounds()
}

func (p *Picture) At(x, y int) color.Color {
	return p.img.At(x, y)
}

func (p *Picture) DrawPixel(x, y int, c color.Color) {
	p.img.Set(x, y, c)
}

func (p *Picture) DrawCircle(cc Circle, w float64, c color.Color) {
	p.DrawWithCondition(
		func(x, y int) bool {
			return float64(x-cc.Center.X)*float64(x-cc.Center.X)+float64(y-cc.Center.Y)*float64(y-cc.Center.Y) <= (cc.Radius+w/2)*(cc.Radius+w/2) && float64(x-cc.Center.X)*float64(x-cc.Center.X)+float64(y-cc.Center.Y)*float64(y-cc.Center.Y) >= (cc.Radius-w/2)*(cc.Radius-w/2)
		},
		[]color.Color{c},
	)
}

func (p *Picture) DrawFilledCircle(cc Circle, c color.Color) {
	p.DrawWithCondition(
		func(x, y int) bool {
			return float64(x-cc.Center.X)*float64(x-cc.Center.X)+float64(y-cc.Center.Y)*float64(y-cc.Center.Y) <= cc.Radius*cc.Radius
		},
		[]color.Color{c},
	)
}

func (p *Picture) DrawRectangle(r image.Rectangle, w float64, c color.Color) {
	p.DrawWithCondition(
		func(x, y int) bool {
			return math.Abs(float64(x-r.Min.X)) <= w || math.Abs(float64(x-r.Max.X)) <= w || math.Abs(float64(y-r.Min.Y)) <= w || math.Abs(float64(y-r.Max.Y)) <= w
		},
		[]color.Color{c},
	)
}

func (p *Picture) DrawFilledRectangle(r image.Rectangle, c color.Color) {
	p.DrawWithCondition(
		func(x, y int) bool {
			return x >= r.Min.X && x < r.Max.X && y >= r.Min.Y && y < r.Max.Y
		},
		[]color.Color{c},
	)
}

func (p *Picture) DrawWithCondition(check func(x, y int) bool, colors []color.Color) error {
	if len(colors) == 0 {
		return fmt.Errorf("drawWithCondition: colors are empty")
	}

	currentColor := 0
	bounds := p.img.Bounds()
	for ix := bounds.Min.X; ix < bounds.Max.X; ix++ {
		for iy := bounds.Min.Y; iy < bounds.Max.Y; iy++ {
			if check(ix, iy) {
				p.img.Set(ix, iy, colors[currentColor%len(colors)])
				currentColor++
			}
		}
	}

	return nil
}

func (p *Picture) SavePNG(name string) error {
	f, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("SavePNG: %w", err)
	}

	err = png.Encode(f, &p.img)
	if err != nil {
		deleteErr := os.Remove(f.Name())
		if deleteErr != nil {
			log.Println(fmt.Errorf("SavePNG: %w", err).Error())
		}

		return fmt.Errorf("SavePNG: %w", err)
	}

	return nil
}
