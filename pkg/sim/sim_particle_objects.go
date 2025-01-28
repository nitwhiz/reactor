package sim

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
	"math/rand/v2"
)

func drawUranium(e *EnvSettings, p *Particle, screen *ebiten.Image) {
	vector.DrawFilledCircle(screen, p.location.X-p.bounds.Radius/2.0, p.location.Y-p.bounds.Radius/2.0, p.bounds.Radius, colornames.Blue, false)
}

func NewUranium() *Particle {
	p := NewParticle()

	p.bounds = &Circle{
		Radius: 10.0,
	}
	p.drawFunc = drawUranium
	p.typ = TypeUranium

	return p
}

func drawNeutron(e *EnvSettings, p *Particle, screen *ebiten.Image) {
	vector.DrawFilledCircle(screen, p.location.X-p.bounds.Radius/2.0, p.location.Y-p.bounds.Radius/2.0, 4, colornames.Grey, false)
}

func reactNeutron(e *Env, self, o Object) {
	if o.Type() == TypeUranium {
		e.Remove(self)
		e.Remove(o)

		for range 3 {
			nE := NewNeutron()

			nE.Location().X = o.Location().X
			nE.Location().Y = o.Location().Y

			nE.Velocity = randVelocity()

			e.Add(nE)
		}
	} else if o.Type() == TypeWater {
		switch w := o.(type) {
		case *Body:
			w.Temperature += e.settings.NeutronWaterHeating

			if w.Temperature < e.settings.WaterEvaporizeTemperature &&
				rand.Float32() < e.settings.WaterNeutronAbsorbRate {
				e.Remove(self)
			}
		}
	}
}

func NewNeutron() *Particle {
	p := NewParticle()

	p.bounds = &Circle{
		Radius: 4.0,
	}
	p.drawFunc = drawNeutron
	p.reactFunc = reactNeutron
	p.typ = TypeNeutron

	return p
}
