package sim

import (
	"github.com/nitwhiz/reactor/pkg/geometry"
	"golang.org/x/image/colornames"
	"image/color"
)

const (
	ParticleElectron = uint64(iota)
	ParticleUranium
	ParticleWater
)

type ParticleType struct {
	Type  uint64
	Color color.Color
}

var ElectronParticleType = &ParticleType{
	Type: ParticleElectron,
	Color: color.RGBA{
		R: 70,
		G: 70,
		B: 70,
		A: 255,
	},
}

var UraniumParticleType = &ParticleType{
	Type:  ParticleUranium,
	Color: colornames.Blue,
}

var NonFissileParticleType = &ParticleType{
	Type:  ParticleUranium,
	Color: colornames.Grey,
}

func CreateElectron(em *EntityManager, x, y, vx, vy float32) {
	em.AddEntity(
		NewParticleTypeComponent(ElectronParticleType),
		NewBodyComponent(geometry.NewCircle(x, y, 5)),
		NewZIndexComponent(100),
		NewVelocityComponent(vx, vy),
	)
}

func CreateUranium(em *EntityManager, x, y float32) {
	em.AddEntity(
		NewParticleTypeComponent(UraniumParticleType),
		NewBodyComponent(geometry.NewCircle(x, y, 10)),
		NewFissionComponent(ParticleElectron),
		NewZIndexComponent(90),
	)
}

func CreateNonFissileElement(em *EntityManager, x, y float32) {
	em.AddEntity(
		NewParticleTypeComponent(NonFissileParticleType),
		NewBodyComponent(geometry.NewCircle(x, y, 10)),
		NewZIndexComponent(80),
	)
}

func CreateWater(em *EntityManager, x, y float32) {
	em.AddEntity(
		NewParticleTypeComponent(&ParticleType{
			Type:  ParticleWater,
			Color: colornames.Lightblue,
		}),
		NewHeatTransferComponent(colornames.Lightblue, ParticleElectron),
		NewBodyComponent(geometry.NewRectangle(x, y, 500, 500)),
		NewZIndexComponent(10),
	)
}
