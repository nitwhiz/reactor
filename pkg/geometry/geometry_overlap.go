package geometry

func circleOverlapsCircle(c1 *Circle, c2 *Circle) bool {
	dx := c1.center.X - c2.center.X
	dy := c1.center.Y - c2.center.Y

	dsq := dx*dx + dy*dy
	radii := c1.Radius + c2.Radius

	return dsq < radii*radii
}

func circleOverlapsRectangle(c *Circle, r *Rectangle) bool {
	rPos := r.TopLeft()

	closestX := max(rPos.X, min(c.center.X, rPos.X+r.Width))
	closestY := max(rPos.Y, min(c.center.Y, rPos.Y+r.Height))

	dx := c.center.X - closestX
	dy := c.center.Y - closestY

	return dx*dx+dy*dy < c.Radius*c.Radius
}

func rectangleOverlapsRectangle(r1 *Rectangle, r2 *Rectangle) bool {
	r1Pos := r1.TopLeft()
	r2Pos := r2.TopLeft()

	return r1Pos.X < r2Pos.X+r2.Width &&
		r1Pos.X+r1.Width > r2Pos.X &&
		r1Pos.Y < r2Pos.Y+r2.Height &&
		r1Pos.Y+r1.Width > r2Pos.Y
}
