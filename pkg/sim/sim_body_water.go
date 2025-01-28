package sim

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
	"time"
)

func updateWater(e *EnvSettings, b *Body, d time.Duration) {
	if !e.UpdateWaterTemperature {
		return
	}

	ds := float32(d) / float32(time.Second)

	if b.Temperature > e.RoomTemperature {
		b.Temperature = max(e.RoomTemperature, b.Temperature-e.WaterTemperatureChangeRate*ds)
	}

	if b.Temperature < e.RoomTemperature {
		b.Temperature = min(e.RoomTemperature, b.Temperature+e.WaterTemperatureChangeRate*ds)
	}
}

func getWaterColor(e *EnvSettings, t float32) color.Color {
	vapeF := float32(math.Pow(float64(t/e.WaterEvaporizeTemperature), 2))

	if vapeF >= 1.0 {
		return color.Transparent
	}

	r := color.RGBA{
		R: 252,
		G: 124,
		B: 97,
		A: 0,
	}

	b := color.RGBA{
		R: 220,
		G: 238,
		B: 255,
		A: 0,
	}

	return color.RGBA{
		R: uint8(float32(b.R) + (float32(r.R)-float32(b.R))*vapeF),
		G: uint8(float32(b.G) + (float32(r.G)-float32(b.G))*vapeF),
		B: uint8(float32(b.B) + (float32(r.B)-float32(b.B))*vapeF),
		A: 255.0,
	}
}

func drawWater(e *EnvSettings, b *Body, screen *ebiten.Image) {
	c := getWaterColor(e, b.Temperature)

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
