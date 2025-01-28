package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nitwhiz/reactor/pkg/sim"
	"golang.org/x/image/colornames"
	"log"
	"math/rand/v2"
)

type Reactor struct {
	env *sim.Env
}

func (r *Reactor) Update() error {
	r.env.Update()

	return nil
}

func (r *Reactor) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.White)

	r.env.Draw(screen)
}

func (r *Reactor) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("reactor")

	reactor := Reactor{
		env: sim.NewEnv(&sim.EnvSettings{
			RoomTemperature:            20.0,
			WaterEvaporizeTemperature:  100.0,
			WaterTemperatureChangeRate: 10.0,
			WaterNeutronAbsorbRate:     0.0125,
			NeutronWaterHeating:        2.0,

			UpdateWaterTemperature: true,
		}),
	}

	for i := range 200 {
		x := i % 15
		y := i / 15

		w := sim.NewWater()

		w.Location().X = float32(x)*25.0 + 150.0
		w.Location().Y = float32(y)*25.0 + 100.0

		reactor.env.Add(w)
	}

	for range 40 {
		e := sim.NewNeutron()

		e.Location().X = 40.0 + 10.0*(rand.Float32()-.5)
		e.Location().Y = 250.0

		e.Velocity.X = 100.0
		e.Velocity.Y = 50.0 * (rand.Float32() - .5)

		reactor.env.Add(e)
	}

	if err := ebiten.RunGame(&reactor); err != nil {
		log.Fatal(err)
	}
}
