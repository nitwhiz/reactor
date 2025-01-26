package sim

import (
	"math"
	"math/rand/v2"
)

func randVelocity() *Velocity {
	a := rand.Float64() * 2.0 * math.Pi

	return &Velocity{
		float32(60.0 * math.Cos(a)),
		float32(60.0 * math.Sin(a)),
	}
}
