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
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("reactor")

	reactor := Reactor{
		env: sim.NewEnv(),
	}

	for i := range 200 {
		x := i % 15
		y := i / 15

		u := sim.NewUranium()

		u.Location().X = float32(x)*25.0 + 150.0
		u.Location().Y = float32(y)*25.0 + 100.0

		reactor.env.Add(u)
	}

	for i := range 10 {
		w := sim.NewWater()

		w.Location().X = 60
		w.Location().Y = 25.0*float32(i) + 40.0

		reactor.env.Add(w)
	}

	e := sim.NewElectron()

	e.Location().X = 40.0
	e.Location().Y = 40.0

	e.Velocity.X = 60.0
	e.Velocity.Y = 60.0

	reactor.env.Add(e)

	if err := ebiten.RunGame(&reactor); err != nil {
		log.Fatal(err)
	}
}
