package sim

import (
	"github.com/hajimehoshi/ebiten/v2"
	"time"
)

type BodyDrawFunc func(e *EnvSettings, b *Body, screen *ebiten.Image)

type BodyUpdateFunc func(e *EnvSettings, b *Body, d time.Duration)

const (
	TypeWater = iota + 1
)

type Body struct {
	id          int64
	drawFunc    BodyDrawFunc
	updateFunc  BodyUpdateFunc
	reactFunc   ReactFunc
	location    *Location
	velocity    *Velocity
	bounds      *Rectangle
	typ         int
	Temperature float32
}

func NewBody() *Body {
	return &Body{
		id:          newId(),
		drawFunc:    nil,
		updateFunc:  nil,
		reactFunc:   nil,
		location:    &Location{},
		velocity:    &Velocity{},
		bounds:      &Rectangle{},
		typ:         TypeNone,
		Temperature: 0.0,
	}
}

func (b *Body) Id() int64 {
	return b.id
}

func (b *Body) Bounds() Bounder {
	return b.bounds
}

func (b *Body) Draw(e *EnvSettings, screen *ebiten.Image) {
	if b.drawFunc != nil {
		b.drawFunc(e, b, screen)
	}
}

func (b *Body) Update(e *EnvSettings, d time.Duration) {
	if b.updateFunc != nil {
		b.updateFunc(e, b, d)
	}
}

func (b *Body) Location() *Location {
	return b.location
}

func (b *Body) Intersects(o Object) bool {
	return b.bounds.Intersects(b.location, o)
}

func (b *Body) React(e *Env, o Object) {
	if b.reactFunc != nil {
		b.reactFunc(e, b, o)
	}
}

func (b *Body) Type() int {
	return b.typ
}
