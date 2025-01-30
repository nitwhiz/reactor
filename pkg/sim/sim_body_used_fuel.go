package sim

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"math/rand/v2"
	"time"
)

type NonUranium struct {
	*Body
}

func NewNonUranium(e *Env) *NonUranium {
	b := NewBody(e)

	b.bounds = &Circle{
		Radius: 8.0,
	}

	return &NonUranium{b}
}

func (n *NonUranium) Draw(screen *ebiten.Image) {
	n.bounds.Draw(n, colornames.Lightgray, screen)
}

func (n *NonUranium) ZIndex() int {
	return 20
}

func (n *NonUranium) Update(_ time.Duration) {
	if rand.Float32() < .0001 {
		u := NewUranium(n.env)

		u.location.X = n.location.X
		u.location.Y = n.location.Y

		n.env.Remove(n)
		n.env.Add(u)
	} else if rand.Float32() < .001 {
		nN := NewNeutron(n.env)

		nN.location.X = n.location.X
		nN.location.Y = n.location.Y

		nN.velocity = randVelocity()

		n.env.Add(nN)
	}
}
