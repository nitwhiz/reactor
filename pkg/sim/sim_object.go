package sim

import (
	"github.com/hajimehoshi/ebiten/v2"
	"time"
)

type Object interface {
	Id() int64
	Bounds() Bounder
	Location() *Location
	Draw(screen *ebiten.Image)
	Update(d time.Duration)
	Intersects(o Object) bool
	React(o Object)
	ZIndex() int
}
