package spatial

type Point2D struct {
	x, y int
}

func NewPoint2D(x, y int) Point2D {
	return Point2D{x, y}
}

func (p Point2D) X() int {
	return p.x
}

func (p Point2D) Y() int {
	return p.y
}

func (p Point2D) Add(v Vector2D) Point2D {
	return Point2D{p.x + v.x, p.y + v.y}
}
