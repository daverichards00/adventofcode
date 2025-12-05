package spatial

var NilBox2D = Box2D{nil: true}

type Box2D struct {
	min Point2D
	max Point2D
	nil bool
}

func NewBox2D(a Point2D, b ...Point2D) Box2D {
	n := Box2D{min: a, max: a}
	for _, p := range b {
		n = n.Expand(p)
	}
	return n
}

func (b Box2D) Min() Point2D {
	return b.min
}

func (b Box2D) Max() Point2D {
	return b.max
}

func (b Box2D) Expand(p Point2D) Box2D {
	if b.nil {
		return Box2D{min: p, max: p}
	}
	return Box2D{
		min: NewPoint2D(min(p.x, b.min.x), min(p.y, b.min.y)),
		max: NewPoint2D(max(p.x, b.max.x), max(p.y, b.max.y)),
	}
}

func (b Box2D) Contains(p Point2D) bool {
	if b.nil {
		return false
	}
	return p.x >= b.min.x && p.x <= b.max.x && p.y >= b.min.y && p.y <= b.max.y
}
