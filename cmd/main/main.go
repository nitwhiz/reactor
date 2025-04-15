package main

import (
	"github.com/nitwhiz/reactor/internal/ebitendisplay"
	"github.com/nitwhiz/reactor/pkg/sim"
	"log"
)

func main() {
	windowWidth := 1280
	windowHeight := 720

	// todo: add control rods
	// todo: add water consuming electrons with a probability

	world := sim.NewEntityManager()

	world.AddSystem(sim.NewMovementSystem(world))
	world.AddSystem(sim.NewHeatTransferSystem(world))
	world.AddSystem(sim.NewFissionSystem(world))
	world.AddSystem(sim.NewWorldBorderSystem(world, -100, -100, float32(windowWidth+100), float32(windowHeight+100)))

	//sim.CreateElectron(world, 15, 500)
	//sim.CreateElectron(world, 15, 520)
	//sim.CreateElectron(world, 15, 540)
	sim.CreateElectron(world, 15, 460, 30, -20)
	sim.CreateElectron(world, 15, 480, 30, -20)
	sim.CreateElectron(world, 15, 500, 30, -20)

	//sim.CreateWater(world, 500, 300)

	//sim.CreateUranium(world, 200, 350)
	//sim.CreateUranium(world, 200, 400)
	//sim.CreateUranium(world, 200, 450)

	for x := float32(0); x < 30; x++ {
		for y := float32(0); y < 20; y++ {
			sim.CreateUranium(world, 100+x*30, 100+y*30)
		}
	}

	r := ebitendisplay.NewReactor(world, windowWidth, windowHeight)

	if err := r.Start(); err != nil {
		log.Fatal(err)
	}
}
