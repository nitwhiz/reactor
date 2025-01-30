package sim

import (
	"sync/atomic"
	"time"
)

var lastId = atomic.Int64{}

func newId() int64 {
	return lastId.Add(1)
}

type ReactFunc func(e *Env, self, o Object)

type Vec2 struct {
	X, Y float32
}

type Location Vec2

type Velocity Vec2

type BodyUpdateFunc func(e *EnvSettings, o Object, d time.Duration)
