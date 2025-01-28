package sim

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
	"time"
)

func updateWater(b *Body, d time.Duration) {
	ds := float32(d) / float32(time.Second)

	if b.Temperature > 20 {
		b.Temperature = max(20.0, b.Temperature-50.0*ds)
	}

	if b.Temperature < 20 {
		b.Temperature = min(20.0, b.Temperature+50.0*ds)
	}
}

func drawWater(b *Body, screen *ebiten.Image) {
	c := colornames.Lightblue

	f := b.Temperature / 20.0

	c.R = uint8(min(255.0, 10.0*f+100.0))
	c.G = uint8(max(120.0, float32(colornames.Lightblue.G)-10.0*f))
	c.B = uint8(max(120.0, float32(colornames.Lightblue.B)-10.0*f))
	c.A = uint8(max(120.0, 255.0-2.0*f))

	vector.DrawFilledRect(screen, b.location.X-b.bounds.Width/2.0, b.location.Y-b.bounds.Height/2.0, b.bounds.Width, b.bounds.Height, c, false)
}

func NewWater() *Body {
	w := NewBody()

	w.bounds = &Rectangle{
		Width:  20,
		Height: 20,
	}
	w.updateFunc = updateWater
	w.drawFunc = drawWater
	w.typ = TypeWater

	return w
}
