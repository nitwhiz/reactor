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
		b.Temperature = min(20, b.Temperature-0.05*ds)
	}

	if b.Temperature < 20 {
		b.Temperature = max(20, b.Temperature+0.05*ds)
	}
}

func drawWater(b *Body, screen *ebiten.Image) {
	c := colornames.Lightblue

	c.R = uint8(min(255.0, 20.0*(b.Temperature/20.0)))

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
