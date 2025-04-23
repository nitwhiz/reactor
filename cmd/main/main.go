package main

import (
	"flag"
	"github.com/nitwhiz/reactor/internal/ebitendisplay"
	"github.com/nitwhiz/reactor/pkg/ecs"
	"github.com/nitwhiz/reactor/pkg/sim"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"strconv"
	"time"
)

func main() {
	flag.BoolFunc("profile", "", func(s string) error {
		secs, _ := strconv.ParseInt(s, 10, 32)

		if secs == 0 {
			secs = 20
		}

		log.Printf("profiling for %d seconds ...", secs)

		cpuProfile, err := os.Create("out/cpu.prof")

		if err != nil {
			log.Fatal(err)
		}

		heapProfile, err := os.Create("out/heap.prof")

		if err != nil {
			log.Fatal(err)
		}

		go func() {
			pprof.StartCPUProfile(cpuProfile)

			time.Sleep(time.Duration(secs) * time.Second)

			pprof.StopCPUProfile()

			pprof.WriteHeapProfile(heapProfile)

			cpuProfile.Close()
			heapProfile.Close()

			os.Exit(0)
		}()

		return nil
	})

	flag.Parse()

	windowWidth := 1280
	windowHeight := 720

	offsetX := float32(160.0)
	offsetY := float32(120.0)

	// todo: add water consuming neutrons with a probability

	em := ecs.NewEntityManager()

	em.AddSystem(sim.NewMovementSystem(em))
	em.AddSystem(sim.NewWaterSystem(em))
	em.AddSystem(sim.NewParticleTemperatureSystem(em))
	em.AddSystem(sim.NewFissionSystem(em))
	em.AddSystem(sim.NewRefillUraniumSystem(em))
	em.AddSystem(sim.NewXenonSystem(em))
	em.AddSystem(sim.NewControlRodSystem(em, 40, offsetY-480/2, offsetY+228, 3.3333))
	em.AddSystem(sim.NewModeratorSystem(em))
	em.AddSystem(sim.NewEmitNeutronsSystem(em))
	em.AddSystem(sim.NewWorldBorderSystem(em, -400, -400, float32(windowWidth+400), float32(windowHeight+400)))

	for x := float32(0); x < 40; x++ {
		for y := float32(0); y < 20; y++ {
			sim.CreateWater(em, offsetX+x*24, offsetY+y*24)

			if rand.Float32() < 0.1 {
				sim.CreateUranium(em, offsetX+x*24, offsetY+y*24)
			} else {
				sim.CreateNonFissileElement(em, offsetX+x*24, offsetY+y*24)
			}
		}
	}

	for x := float32(0); x < 10; x++ {
		if int(x)%2 == 0 {
			sim.CreateMovableControlRod(em, offsetX+x*24*4+36, offsetY+228, sim.TagControlRodSet1)
		} else {
			sim.CreateStaticControlRod(em, offsetX+x*24*4+36, offsetY+228, sim.TagControlRodSet2)
		}
	}

	for x := float32(0); x < 11; x++ {
		sim.CreateModerator(em, offsetX+x*24*4-12, offsetY+228)
	}

	r := ebitendisplay.NewReactor(em, windowWidth, windowHeight)

	if err := r.Start(); err != nil {
		log.Fatal(err)
	}
}
