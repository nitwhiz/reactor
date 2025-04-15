package sim

import (
	"math"
	"math/rand"
)

type MovementSystem struct {
	world     *EntityManager
	signature Signature
}

func NewMovementSystem(world *EntityManager) *MovementSystem {
	return &MovementSystem{
		world:     world,
		signature: ComponentTypeBody | ComponentTypeVelocity,
	}
}

func (s *MovementSystem) Update() {
	s.world.EachEntity(s.signature, func(q *Query, eId EntityID) bool {
		b := s.world.GetComponent(ComponentTypeBody, eId).(*BodyComponent)
		v := s.world.GetComponent(ComponentTypeVelocity, eId).(*VelocityComponent)

		loc := b.Body.Location()

		loc.X += v.velocity.X / 16.6
		loc.Y += v.velocity.Y / 16.6

		return true
	})
}

type HeatTransferSystem struct {
	world     *EntityManager
	signature Signature
}

func NewHeatTransferSystem(world *EntityManager) *HeatTransferSystem {
	return &HeatTransferSystem{
		world:     world,
		signature: ComponentTypeParticleType | ComponentTypeHeatTransfer,
	}
}

func (s *HeatTransferSystem) Update() {
	s.world.EachEntity(s.signature, func(q *Query, eId EntityID) bool {
		p := s.world.GetComponent(ComponentTypeParticleType, eId).(*ParticleTypeComponent)
		h := s.world.GetComponent(ComponentTypeHeatTransfer, eId).(*HeatTransferComponent)

		hadCollision := false

		eachCollision(s.world, eId, ComponentTypeParticleType, func(otherEntityId EntityID, _ *BodyComponent) bool {
			otherParticleTypeComponent := s.world.GetComponent(ComponentTypeParticleType, otherEntityId).(*ParticleTypeComponent)

			if otherParticleTypeComponent.ParticleType.Type == h.ParticleType {
				hadCollision = true
				h.Temperature = min(100.0, h.Temperature+1.0/16.66)
			}

			return true
		})

		if !hadCollision {
			h.Temperature = max(20.0, h.Temperature-1.0/16.66)
		}

		p.ParticleType.Color = h.CurrentColor()

		return true
	})
}

type WorldBorderSystem struct {
	world *EntityManager
	minX  float32
	minY  float32
	maxX  float32
	maxY  float32
}

func NewWorldBorderSystem(world *EntityManager, minX, minY, maxX, maxY float32) *WorldBorderSystem {
	return &WorldBorderSystem{
		world: world,
		minX:  minX,
		minY:  minY,
		maxX:  maxX,
		maxY:  maxY,
	}
}

func (s *WorldBorderSystem) Update() {
	s.world.EachEntity(ComponentTypeBody, func(q *Query, eId EntityID) bool {
		b := s.world.GetComponent(ComponentTypeBody, eId).(*BodyComponent)
		loc := b.Body.Location()

		if loc.X < s.minX || loc.Y < s.minY || loc.X > s.maxX || loc.Y > s.maxY {
			s.world.RemoveEntity(eId)
		}

		return true
	})
}

type FissionSystem struct {
	world     *EntityManager
	signature Signature
}

func NewFissionSystem(world *EntityManager) *FissionSystem {
	return &FissionSystem{
		world:     world,
		signature: ComponentTypeFission,
	}
}

func (s *FissionSystem) Update() {
	v := 30.0

	s.world.EachEntity(s.signature, func(q *Query, eId EntityID) bool {
		f := s.world.GetComponent(ComponentTypeFission, eId).(*FissionComponent)

		eachCollision(s.world, eId, ComponentTypeParticleType, func(otherEntityId EntityID, b *BodyComponent) bool {
			otherParticleTypeComponent := s.world.GetComponent(ComponentTypeParticleType, otherEntityId).(*ParticleTypeComponent)

			if otherParticleTypeComponent.ParticleType.Type == f.InducingParticleType {
				loc := b.Body.Location()

				CreateNonFissileElement(s.world, loc.X, loc.Y)

				for range 3 {
					a := math.Pi * 2 * rand.Float64()

					vx := v * math.Cos(a)
					vy := v * math.Sin(a)

					CreateElectron(s.world, loc.X, loc.Y, float32(vx), float32(vy))
				}

				s.world.RemoveEntity(eId)

				return false
			}

			return true
		})

		return true
	})
}
