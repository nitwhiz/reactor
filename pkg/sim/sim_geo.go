package sim

func clamp(val, min, max float32) float32 {
	if val < min {
		return min
	}

	if val > max {
		return max
	}

	return val
}

func distanceSquared(x1, y1, x2, y2 float32) float32 {
	dx := x1 - x2
	dy := y1 - y2

	return dx*dx + dy*dy
}

func circleIntersectsRectangle(l *Location, c *Circle, oLoc *Location, b *Rectangle) bool {
	closestX := clamp(l.X, oLoc.X-b.Width/2.0, oLoc.X+b.Width/2.0)
	closestY := clamp(l.Y, oLoc.Y-b.Height/2.0, oLoc.Y+b.Height/2.0)

	return distanceSquared(closestX, closestY, l.X, l.Y) <= c.Radius
}

type Bounder interface {
	Intersects(l *Location, o Object) bool
}

type Circle struct {
	Radius float32
}

func (c *Circle) Intersects(l *Location, o Object) bool {
	oLoc := o.Location()

	switch b := o.Bounds().(type) {
	case *Circle:
		ds := distanceSquared(l.X, l.Y, oLoc.X, oLoc.Y)
		dr := c.Radius - b.Radius
		sr := c.Radius + b.Radius

		return dr*dr <= ds && ds <= sr*sr
	case *Rectangle:
		return circleIntersectsRectangle(l, c, oLoc, b)
	}

	return false
}

type Rectangle struct {
	Width, Height float32
}

func (r *Rectangle) Intersects(l *Location, o Object) bool {
	oLoc := o.Location()

	switch b := o.Bounds().(type) {
	case *Circle:
		return circleIntersectsRectangle(oLoc, b, oLoc, r)
	case *Rectangle:
		rMinX := l.X - r.Width/2.0
		rMaxX := l.X + r.Width/2.0
		rMinY := l.Y - r.Height/2.0
		rMaxY := l.Y + r.Height/2.0

		oMinX := oLoc.X - b.Width/2.0
		oMaxX := oLoc.X + b.Width/2.0
		oMinY := oLoc.Y - b.Height/2.0
		oMaxY := oLoc.Y + b.Height/2.0

		return !(rMaxX < oMinX || rMinX > oMaxX || rMaxY < oMinY || rMinY > oMaxY)
	}

	return false
}
