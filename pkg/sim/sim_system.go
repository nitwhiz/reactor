package sim

import (
	"github.com/nitwhiz/reactor/pkg/ecs"
	"math"
	"math/rand"
)

type MovementSystem struct {
	em *ecs.EntityManager
}

func NewMovementSystem(em *ecs.EntityManager) *MovementSystem {
	return &MovementSystem{
		em: em,
	}
}

func (s *MovementSystem) Update() {
	s.em.EachEntity(ComponentTypeBody|ComponentTypeVelocity, func(q *ecs.Query, eId ecs.EntityID) bool {
		b := s.em.GetComponent(ComponentTypeBody, eId).(*BodyComponent)
		v := s.em.GetComponent(ComponentTypeVelocity, eId).(*VelocityComponent)

		loc := b.Body.Location()

		loc.X += v.velocity.X / 16.6
		loc.Y += v.velocity.Y / 16.6

		return true
	})
}

type WaterSystem struct {
	em *ecs.EntityManager
}

func NewWaterSystem(em *ecs.EntityManager) *WaterSystem {
	return &WaterSystem{
		em: em,
	}
}

func (s *WaterSystem) Update() {
	s.em.EachEntity(TagWater|ComponentTypeTemperature|ComponentTypeBody, func(q *ecs.Query, eId ecs.EntityID) bool {
		t := s.em.GetComponent(ComponentTypeTemperature, eId).(*TemperatureComponent)
		b := s.em.GetComponent(ComponentTypeBody, eId).(*BodyComponent)

		hadCollision := false

		eachBodyCollision(s.em, TagThermalNeutron, b, func() bool {
			hadCollision = true
			t.Temperature = t.Temperature + 6.0/16.66

			return true
		})

		// todo: how accurate is this?
		eachBodyCollision(s.em, TagFastNeutron, b, func() bool {
			hadCollision = true
			t.Temperature = t.Temperature + 6.0/16.66

			return true
		})

		if !hadCollision {
			t.Temperature = max(20.0, t.Temperature-6.0/16.66)
		}

		if t.Temperature < 100 {
			eachCollision(s.em, eId, TagThermalNeutron, func(neutronEntityId ecs.EntityID, _ *BodyComponent) bool {
				if rand.Float32() < .002/(t.Temperature-19.0) {
					s.em.RemoveEntity(neutronEntityId)
				}

				return true
			})

			eachCollision(s.em, eId, TagFastNeutron, func(neutronEntityId ecs.EntityID, _ *BodyComponent) bool {
				if rand.Float32() < .002/(t.Temperature-19.0) {
					s.em.RemoveEntity(neutronEntityId)
				}

				return true
			})
		}

		return true
	})
}

type ParticleTemperatureSystem struct {
	em *ecs.EntityManager
}

func NewParticleTemperatureSystem(em *ecs.EntityManager) *ParticleTemperatureSystem {
	return &ParticleTemperatureSystem{
		em: em,
	}
}

func (s *ParticleTemperatureSystem) Update() {
	s.em.EachEntity(ComponentTypeTemperature|ComponentTypeParticle, func(q *ecs.Query, eId ecs.EntityID) bool {
		t := s.em.GetComponent(ComponentTypeTemperature, eId).(*TemperatureComponent)
		p := s.em.GetComponent(ComponentTypeParticle, eId).(*ParticleComponent)

		p.ParticleType.Color = t.CurrentColor()

		return true
	})
}

type WorldBorderSystem struct {
	em   *ecs.EntityManager
	minX float32
	minY float32
	maxX float32
	maxY float32
}

func NewWorldBorderSystem(em *ecs.EntityManager, minX, minY, maxX, maxY float32) *WorldBorderSystem {
	return &WorldBorderSystem{
		em:   em,
		minX: minX,
		minY: minY,
		maxX: maxX,
		maxY: maxY,
	}
}

func (s *WorldBorderSystem) Update() {
	s.em.EachEntity(ComponentTypeBody|ComponentTypeVelocity, func(q *ecs.Query, eId ecs.EntityID) bool {
		b := s.em.GetComponent(ComponentTypeBody, eId).(*BodyComponent)
		loc := b.Body.Location()

		if loc.X < s.minX || loc.Y < s.minY || loc.X > s.maxX || loc.Y > s.maxY {
			s.em.RemoveEntity(eId)
		}

		return true
	})
}

type FissionSystem struct {
	em *ecs.EntityManager
}

func NewFissionSystem(em *ecs.EntityManager) *FissionSystem {
	return &FissionSystem{
		em: em,
	}
}

func (s *FissionSystem) Update() {
	s.em.EachEntity(ComponentTypeBody|TagFission, func(q *ecs.Query, eId ecs.EntityID) bool {
		b := s.em.GetComponent(ComponentTypeBody, eId).(*BodyComponent)

		eachBodyCollision(s.em, TagThermalNeutron, b, func() bool {
			// fixme: this collision loop removes entities multiple times (which, why?)
			// as soon as a collision is found, this loop returns false, and other electrons will not be checked

			loc := b.Body.Location()

			CreateNonFissileElement(s.em, loc.X, loc.Y)

			for range 3 {
				CreateFastNeutron(s.em, loc.X, loc.Y, math.Pi*2*rand.Float64())
			}

			s.em.RemoveEntity(eId)

			return false
		})

		return true
	})
}

type ControlRodSystem struct {
	em                 *ecs.EntityManager
	targetNeutronCount int
	minY               float32
	maxY               float32
	speed              float32
	moveSet            int
}

func NewControlRodSystem(em *ecs.EntityManager, targetNeutronCount int, minY float32, maxY float32, speed float32) *ControlRodSystem {
	return &ControlRodSystem{
		em:                 em,
		targetNeutronCount: targetNeutronCount,
		minY:               minY,
		maxY:               maxY,
		speed:              speed,
		moveSet:            0,
	}
}

func (s *ControlRodSystem) Update() {
	s.em.EachEntity(TagControlRod, func(q *ecs.Query, eId ecs.EntityID) bool {
		eachCollision(s.em, eId, TagThermalNeutron, func(otherEntityId ecs.EntityID, b *BodyComponent) bool {
			s.em.RemoveEntity(otherEntityId)

			return true
		})

		return true
	})

	neutronCount := s.em.CountComponent(TagThermalNeutron) + s.em.CountComponent(TagFastNeutron)

	if neutronCount > s.targetNeutronCount {
		if s.moveSet == 0 {
			s.em.EachEntity(TagControlRodSet1|ComponentTypeVelocity|ComponentTypeBody, func(q *ecs.Query, eId ecs.EntityID) bool {
				v := s.em.GetComponent(ComponentTypeVelocity, eId).(*VelocityComponent)
				b := s.em.GetComponent(ComponentTypeBody, eId).(*BodyComponent)

				if b.Body.Location().Y > s.maxY {
					s.moveSet = 1
					v.velocity.Y = 0
				} else {
					v.velocity.Y = s.speed
				}

				return true
			})
		} else {
			s.em.EachEntity(TagControlRodSet2|ComponentTypeVelocity|ComponentTypeBody, func(q *ecs.Query, eId ecs.EntityID) bool {
				v := s.em.GetComponent(ComponentTypeVelocity, eId).(*VelocityComponent)
				b := s.em.GetComponent(ComponentTypeBody, eId).(*BodyComponent)

				if b.Body.Location().Y > s.maxY {
					s.moveSet = 0
					v.velocity.Y = 0
				} else {
					v.velocity.Y = s.speed
				}

				return true
			})
		}
	} else {
		if s.moveSet == 0 {
			s.em.EachEntity(TagControlRodSet1|ComponentTypeVelocity, func(q *ecs.Query, eId ecs.EntityID) bool {
				v := s.em.GetComponent(ComponentTypeVelocity, eId).(*VelocityComponent)
				b := s.em.GetComponent(ComponentTypeBody, eId).(*BodyComponent)

				if b.Body.Location().Y < s.minY {
					s.moveSet = 1
					v.velocity.Y = 0
				} else {
					v.velocity.Y = -s.speed
				}

				return true
			})
		} else {
			s.em.EachEntity(TagControlRodSet2|ComponentTypeVelocity, func(q *ecs.Query, eId ecs.EntityID) bool {
				v := s.em.GetComponent(ComponentTypeVelocity, eId).(*VelocityComponent)
				b := s.em.GetComponent(ComponentTypeBody, eId).(*BodyComponent)

				if b.Body.Location().Y < s.minY {
					s.moveSet = 0
					v.velocity.Y = 0
				} else {
					v.velocity.Y = -s.speed
				}

				return true
			})
		}
	}
}

type EmitNeutronsSystem struct {
	em *ecs.EntityManager
}

func NewEmitNeutronsSystem(em *ecs.EntityManager) *EmitNeutronsSystem {
	return &EmitNeutronsSystem{
		em: em,
	}
}

func (s *EmitNeutronsSystem) Update() {
	s.em.EachEntity(TagEmitNeutrons|ComponentTypeBody, func(q *ecs.Query, eId ecs.EntityID) bool {
		if rand.Float64() < .00004 {
			b := s.em.GetComponent(ComponentTypeBody, eId).(*BodyComponent)
			loc := b.Body.Location()

			CreateFastNeutron(s.em, loc.X, loc.Y, math.Pi*2*rand.Float64())
		}

		return true
	})

}

type RefillUraniumSystem struct {
	em *ecs.EntityManager
}

func NewRefillUraniumSystem(em *ecs.EntityManager) *RefillUraniumSystem {
	return &RefillUraniumSystem{
		em: em,
	}
}

func (s *RefillUraniumSystem) Update() {
	s.em.EachEntity(ComponentTypeBody|TagNonFissile, func(q *ecs.Query, eId ecs.EntityID) bool {
		if rand.Float32() < .00006 {
			s.em.RemoveEntity(eId)

			b := s.em.GetComponent(ComponentTypeBody, eId).(*BodyComponent)

			bLoc := b.Body.Location()

			CreateUranium(s.em, bLoc.X, bLoc.Y)
		}

		return true
	})
}

type XenonSystem struct {
	em *ecs.EntityManager
}

func NewXenonSystem(em *ecs.EntityManager) *XenonSystem {
	return &XenonSystem{
		em: em,
	}
}

func (s *XenonSystem) Update() {
	s.em.EachEntity(ComponentTypeBody|TagNonFissile, func(q *ecs.Query, eId ecs.EntityID) bool {
		if rand.Float32() < .00002 {
			s.em.RemoveEntity(eId)

			b := s.em.GetComponent(ComponentTypeBody, eId).(*BodyComponent)

			bLoc := b.Body.Location()

			CreateXenon(s.em, bLoc.X, bLoc.Y)
		}

		return true
	})

	s.em.EachEntity(TagXenon, func(q *ecs.Query, eId ecs.EntityID) bool {
		eachCollision(s.em, eId, TagThermalNeutron, func(otherEntityId ecs.EntityID, b *BodyComponent) bool {
			s.em.RemoveEntity(otherEntityId)
			s.em.RemoveEntity(eId)

			bLoc := b.Body.Location()

			CreateNonFissileElement(s.em, bLoc.X, bLoc.Y)

			return false
		})

		return true
	})
}

type ModeratorSystem struct {
	em *ecs.EntityManager
}

func NewModeratorSystem(em *ecs.EntityManager) *ModeratorSystem {
	return &ModeratorSystem{
		em: em,
	}
}

func (s *ModeratorSystem) Update() {
	s.em.EachEntity(TagModerator, func(q *ecs.Query, eId ecs.EntityID) bool {
		eachCollision(s.em, eId, TagFastNeutron|ComponentTypeVelocity, func(otherEntityId ecs.EntityID, b *BodyComponent) bool {
			s.em.RemoveEntity(otherEntityId)

			v := s.em.GetComponent(ComponentTypeVelocity, otherEntityId).(*VelocityComponent)
			angle := math.Atan2(float64(v.velocity.Y), float64(v.velocity.X))

			bLoc := b.Body.Location()

			CreateThermalNeutron(s.em, bLoc.X, bLoc.Y, math.Pi-angle)

			return true
		})

		return true
	})
}
