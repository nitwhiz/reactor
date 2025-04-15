package geometry

type Vec2 struct {
	X float32
	Y float32
}

type Body interface {
	Location() *Vec2
	Overlaps(l Body) bool
}

type Circle struct {
	Radius float32
	center *Vec2
}

func NewCircle(x, y float32, radius float32) *Circle {
	return &Circle{
		center: &Vec2{
			X: x,
			Y: y,
		},
		Radius: radius,
	}
}

func (c *Circle) Location() *Vec2 {
	return c.center
}

func (c *Circle) Overlaps(other Body) bool {
	switch o := other.(type) {
	case *Rectangle:
		return circleOverlapsRectangle(c, o)
	case *Circle:
		return circleOverlapsCircle(c, o)
	}

	return false
}

type Rectangle struct {
	Width  float32
	Height float32
	center *Vec2
}

func NewRectangle(x, y float32, width, height float32) *Rectangle {
	return &Rectangle{
		Width:  width,
		Height: height,
		center: &Vec2{
			X: x,
			Y: y,
		},
	}
}

func (r *Rectangle) Location() *Vec2 {
	return r.center
}

func (r *Rectangle) TopLeft() Vec2 {
	return Vec2{
		X: r.center.X - r.Width/2.0,
		Y: r.center.Y - r.Height/2.0,
	}
}

func (r *Rectangle) Overlaps(other Body) bool {
	switch o := other.(type) {
	case *Rectangle:
		return rectangleOverlapsRectangle(r, o)
	case *Circle:
		return circleOverlapsRectangle(o, r)
	}

	return false
}
