package sim

import (
	"github.com/hajimehoshi/ebiten/v2"
	"sync/atomic"
	"time"
)

var lastId = atomic.Int64{}

func newId() int64 {
	return lastId.Add(1)
}

type ReactFunc func(e *Env, self, o Object)

type Object interface {
	Id() int64
	Bounds() Bounder
	Draw(screen *ebiten.Image)
	Update(d time.Duration)
	Location() *Location
	Intersects(o Object) bool
	React(e *Env, o Object)
	Type() int
}

type Vec2 struct {
	X, Y float32
}

type Location Vec2

type Velocity Vec2
