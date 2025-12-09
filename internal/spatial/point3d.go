package spatial

type Point3D struct {
	x, y, z int
}

func NewPoint3D(x, y, z int) Point3D {
	return Point3D{x, y, z}
}

func (p Point3D) X() int {
	return p.x
}

func (p Point3D) Y() int {
	return p.y
}

func (p Point3D) Z() int {
	return p.z
}
