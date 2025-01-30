package sim

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

type Uranium struct {
	*Body
}

func NewUranium(e *Env) *Uranium {
	b := NewBody(e)

	b.bounds = &Circle{
		Radius: 8.0,
	}

	return &Uranium{b}
}

func (u *Uranium) Draw(screen *ebiten.Image) {
	u.bounds.Draw(u, colornames.Blue, screen)
}

func (u *Uranium) ZIndex() int {
	return 20
}
