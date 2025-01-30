package sim

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

type ControlRod struct {
	*Body
}

func NewControlRod(e *Env) *ControlRod {
	b := NewBody(e)

	b.bounds = &Rectangle{
		Width:  5,
		Height: 400,
	}

	return &ControlRod{b}
}

func (c *ControlRod) Draw(screen *ebiten.Image) {
	c.bounds.Draw(c, colornames.Darkgray, screen)
}

func (c *ControlRod) ZIndex() int {
	return 30
}

func (c *ControlRod) React(o Object) {
	switch o.(type) {
	case *Neutron:
		c.env.Remove(o)
	}
}
