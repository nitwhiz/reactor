package sim

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"math"
	"time"
)

type Water struct {
	*Body
	Temperature float32
}

func NewWater(e *Env) *Water {
	w := NewBody(e)

	w.bounds = &Rectangle{
		Width:  20,
		Height: 20,
	}

	return &Water{Body: w, Temperature: e.settings.RoomTemperature}
}

func (w *Water) ZIndex() int {
	return 10
}

func (w *Water) Update(d time.Duration) {
	if !w.env.settings.UpdateWaterTemperature {
		return
	}

	ds := float32(d) / float32(time.Second)

	if w.Temperature > w.env.settings.RoomTemperature {
		w.Temperature = max(w.env.settings.RoomTemperature, w.Temperature-w.env.settings.WaterTemperatureChangeRate*ds)
	}

	if w.Temperature < w.env.settings.RoomTemperature {
		w.Temperature = min(w.env.settings.RoomTemperature, w.Temperature+w.env.settings.WaterTemperatureChangeRate*ds)
	}
}

func (w *Water) Draw(screen *ebiten.Image) {
	w.bounds.Draw(w, w.getColor(), screen)
}

func (w *Water) getColor() color.Color {
	vapeF := float32(math.Pow(float64(w.Temperature/w.env.settings.WaterVaporizeTemperature), 2))

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
