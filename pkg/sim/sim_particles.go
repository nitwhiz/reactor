package sim

import (
	"github.com/nitwhiz/reactor/pkg/ecs"
	"github.com/nitwhiz/reactor/pkg/geometry"
	"image/color"
	"math"
)

type ParticleType struct {
	Color color.Color
}

var thermalNeutronParticleType = &ParticleType{
	Color: color.RGBA{
		R: 90,
		G: 90,
		B: 90,
		A: 255,
	},
}

var uraniumParticleType = &ParticleType{
	Color: color.RGBA{
		R: 37,
		G: 139,
		B: 254,
		A: 255,
	},
}

var nonFissileParticleType = &ParticleType{
	Color: color.RGBA{
		R: 190,
		G: 190,
		B: 190,
		A: 255,
	},
}

var controlRodParticleType = &ParticleType{
	Color: color.RGBA{
		R: 70,
		G: 70,
		B: 70,
		A: 255,
	},
}

var xenonParticleType = &ParticleType{
	Color: color.RGBA{
		R: 70,
		G: 70,
		B: 70,
		A: 255,
	},
}

var fastNeutronParticleType = &ParticleType{
	Color: color.RGBA{
		R: 120,
		G: 120,
		B: 120,
		A: 255,
	},
}

var moderatorParticleType = &ParticleType{
	Color: color.RGBA{
		R: 190,
		G: 190,
		B: 190,
		A: 255,
	},
}

func CreateThermalNeutron(em *ecs.EntityManager, x, y float32, angle float64) {
	vx := float32(30.0 * math.Cos(angle))
	vy := float32(30.0 * math.Sin(angle))

	em.AddEntity(
		ecs.NewTagComponent(TagThermalNeutron),
		NewParticleTypeComponent(thermalNeutronParticleType),
		NewBodyComponent(geometry.NewCircle(x, y, 5)),
		NewRenderComponent(100),
		NewVelocityComponent(vx, vy),
	)
}

func CreateUranium(em *ecs.EntityManager, x, y float32) {
	em.AddEntity(
		ecs.NewTagComponent(TagFission),
		NewParticleTypeComponent(uraniumParticleType),
		NewBodyComponent(geometry.NewCircle(x, y, 8)),
		NewRenderComponent(80),
	)
}

func CreateNonFissileElement(em *ecs.EntityManager, x, y float32) {
	em.AddEntity(
		NewParticleTypeComponent(nonFissileParticleType),
		NewBodyComponent(geometry.NewCircle(x, y, 8)),
		NewRenderComponent(70),
		ecs.NewTagComponent(TagEmitNeutrons),
		ecs.NewTagComponent(TagNonFissile),
	)
}

func CreateWater(em *ecs.EntityManager, x, y float32) {
	em.AddEntity(
		ecs.NewTagComponent(TagWater),
		NewParticleTypeComponent(&ParticleType{Color: color.Transparent}),
		NewTemperatureComponent(color.RGBA{
			R: 224,
			G: 237,
			B: 249,
			A: 255,
		}),
		NewBodyComponent(geometry.NewRectangle(x, y, 22, 22)),
		NewRenderComponent(10),
	)
}

func CreateMovableControlRod(em *ecs.EntityManager, x, y float32, tag ecs.ComponentType) {
	em.AddEntity(
		ecs.NewTagComponent(tag),
		NewParticleTypeComponent(controlRodParticleType),
		NewBodyComponent(geometry.NewRectangle(x, y, 6, 480)),
		NewVelocityComponent(0, 0),
		NewRenderComponent(90),
	)
}

func CreateStaticControlRod(em *ecs.EntityManager, x, y float32, tag ecs.ComponentType) {
	em.AddEntity(
		ecs.NewTagComponent(tag),
		ecs.NewTagComponent(TagControlRod),
		NewParticleTypeComponent(controlRodParticleType),
		NewBodyComponent(geometry.NewRectangle(x, y, 6, 480)),
		NewVelocityComponent(0, 0),
		NewRenderComponent(90),
	)
}

func CreateModerator(em *ecs.EntityManager, x, y float32) {
	em.AddEntity(
		ecs.NewTagComponent(TagModerator),
		NewParticleTypeComponent(moderatorParticleType),
		NewBodyComponent(geometry.NewRectangle(x, y, 6, 480)),
		NewVelocityComponent(0, 0),
		NewRenderComponent(90),
	)
}

func CreateXenon(em *ecs.EntityManager, x, y float32) {
	em.AddEntity(
		ecs.NewTagComponent(TagXenon),
		NewParticleTypeComponent(xenonParticleType),
		NewBodyComponent(geometry.NewCircle(x, y, 8)),
		NewRenderComponent(60),
	)
}

func CreateFastNeutron(em *ecs.EntityManager, x, y float32, angle float64) {
	vx := float32(60.0 * math.Cos(angle))
	vy := float32(60.0 * math.Sin(angle))

	em.AddEntity(
		ecs.NewTagComponent(TagFastNeutron),
		NewParticleTypeComponent(fastNeutronParticleType),
		NewBodyComponent(geometry.NewCircle(x, y, 5)),
		NewVelocityComponent(vx, vy),
		NewRenderComponent(110),
	)
}
