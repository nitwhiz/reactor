package sim

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
	"time"
)

type ParticleDrawFunc func(e *EnvSettings, p *Particle, screen *ebiten.Image)

type ParticleUpdateFunc func(e *EnvSettings, p *Particle, d time.Duration)

var defaultParticleDraw = func(e *EnvSettings, p *Particle, screen *ebiten.Image) {
	vector.DrawFilledCircle(screen, p.location.X-1, p.location.Y-1, 3, colornames.White, false)
}

const (
	TypeNone = iota
	TypeNeutron
	TypeUranium
)

type Particle struct {
	id         int64
	drawFunc   ParticleDrawFunc
	updateFunc ParticleUpdateFunc
	reactFunc  ReactFunc
	location   *Location
	bounds     *Circle
	typ        int
	Velocity   *Velocity
}

func NewParticle() *Particle {
	return &Particle{
		id:         newId(),
		drawFunc:   defaultParticleDraw,
		updateFunc: nil,
		reactFunc:  nil,
		location:   &Location{},
		bounds:     &Circle{},
		typ:        TypeNone,
		Velocity:   &Velocity{},
	}
}

func (p *Particle) Id() int64 {
	return p.id
}

func (p *Particle) Type() int {
	return p.typ
}

func (p *Particle) Bounds() Bounder {
	return p.bounds
}

func (p *Particle) Location() *Location {
	return p.location
}

func (p *Particle) Update(e *EnvSettings, d time.Duration) {
	ds := float32(d) / float32(time.Second)

	p.location.X += p.Velocity.X * ds
	p.location.Y += p.Velocity.Y * ds

	if p.updateFunc != nil {
		p.updateFunc(e, p, d)
	}
}

func (p *Particle) Draw(e *EnvSettings, screen *ebiten.Image) {
	if p.drawFunc != nil {
		p.drawFunc(e, p, screen)
	}
}

func (p *Particle) Intersects(o Object) bool {
	return p.bounds.Intersects(p.location, o)
}

func (p *Particle) React(e *Env, o Object) {
	if p.reactFunc != nil {
		p.reactFunc(e, p, o)
	}
}
