package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nitwhiz/reactor/pkg/sim"
	"golang.org/x/image/colornames"
	"log"
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
	return 1280, 720
}

func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("reactor")

	reactor := Reactor{
		env: sim.NewEnv(&sim.EnvSettings{
			RoomTemperature:            20.0,
			WaterVaporizeTemperature:   100.0,
			WaterTemperatureChangeRate: 10.0,
			WaterNeutronAbsorbRate:     0.0125,
			NeutronWaterHeating:        2.0,

			UpdateWaterTemperature: true,
		}),
	}

	for y := range 29 {
		for x := range 54 {
			w := sim.NewWater(reactor.env)

			wLoc := w.Location()

			wLoc.X = float32(x)*22.0 + 50.0
			wLoc.Y = float32(y)*22.0 + 50.0

			reactor.env.Add(w)

			u := sim.NewNonUranium(reactor.env)

			uLoc := u.Location()

			uLoc.X = wLoc.X
			uLoc.Y = wLoc.Y

			reactor.env.Add(u)
		}
	}

	//for i := range 200 {
	//	x := i % 15
	//	y := i / 15
	//
	//	w := sim.NewWater(reactor.env)
	//
	//	w.Location().X = float32(x)*25.0 + 150.0
	//	w.Location().Y = float32(y)*25.0 + 100.0
	//
	//	reactor.env.Add(w)
	//
	//	u := sim.NewNonUranium(reactor.env)
	//
	//	u.Location().X = w.Location().X
	//	u.Location().Y = w.Location().Y
	//
	//	reactor.env.Add(u)
	//}
	//
	//c := sim.NewControlRod(reactor.env)
	//
	//c.Location().X = 600
	//c.Location().Y = 250
	//
	//reactor.env.Add(c)

	//for range 40 {
	//	e := sim.NewNeutron(reactor.env)
	//
	//	e.Location().X = 40.0 + 10.0*(rand.Float32()-.5)
	//	e.Location().Y = 250.0
	//
	//	e.Velocity().X = 100.0
	//	e.Velocity().Y = 50.0 * (rand.Float32() - .5)
	//
	//	reactor.env.Add(e)
	//}

	if err := ebiten.RunGame(&reactor); err != nil {
		log.Fatal(err)
	}
}
