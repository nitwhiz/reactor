package sim

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"image/color"
	"time"
)

type Body struct {
	id       int64
	location *Location
	velocity *Velocity
	bounds   Bounder
	env      *Env
}

func NewBody(env *Env) *Body {
	return &Body{
		id:       newId(),
		location: &Location{},
		velocity: &Velocity{},
		bounds:   &Rectangle{},
		env:      env,
	}
}

func (b *Body) Color() color.Color {
	return colornames.Red
}

func (b *Body) ZIndex() int {
	return 10
}

func (b *Body) Id() int64 {
	return b.id
}

func (b *Body) Bounds() Bounder {
	return b.bounds
}

func (b *Body) Draw(screen *ebiten.Image) {
	b.bounds.Draw(b, colornames.Red, screen)
}

func (b *Body) Update(d time.Duration) {
	ds := float32(d) / float32(time.Second)

	b.location.X += b.velocity.X * ds
	b.location.Y += b.velocity.Y * ds
}

func (b *Body) Location() *Location {
	return b.location
}

func (b *Body) Velocity() *Velocity {
	return b.velocity
}

func (b *Body) Intersects(o Object) bool {
	return b.bounds.Intersects(b.location, o)
}

func (b *Body) React(o Object) {
}
