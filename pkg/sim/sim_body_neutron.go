package sim

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"math/rand/v2"
)

type Neutron struct {
	*Body
}

func NewNeutron(e *Env) *Neutron {
	b := NewBody(e)

	b.bounds = &Circle{
		Radius: 4.0,
	}

	return &Neutron{b}
}

func (n *Neutron) ZIndex() int {
	return 1000
}

func (n *Neutron) React(other Object) {
	switch o := other.(type) {
	case *Uranium:
		n.env.Remove(n)
		n.env.Remove(o)

		for range 3 {
			nE := NewNeutron(n.env)

			nE.Location().X = o.Location().X
			nE.Location().Y = o.Location().Y

			nE.velocity = randVelocity()

			n.env.Add(nE)
		}

		nN := NewNonUranium(n.env)

		nN.location.X = o.location.X
		nN.location.Y = o.location.Y

		n.env.Add(nN)
	case *Water:
		o.Temperature += n.env.settings.NeutronWaterHeating

		if o.Temperature < n.env.settings.WaterVaporizeTemperature &&
			rand.Float32() < n.env.settings.WaterNeutronAbsorbRate {
			n.env.Remove(n)
		}
	}
}

func (n *Neutron) Draw(screen *ebiten.Image) {
	n.bounds.Draw(n, colornames.Grey, screen)
}
