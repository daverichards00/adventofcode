package spatial

import "fmt"

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

func (p Point3D) Add(v Vector3D) Point3D {
	return Point3D{p.x + v.x, p.y + v.y, p.z + v.z}
}

func (p Point3D) To(other Point3D) Vector3D {
	return Vector3D{other.x - p.x, other.y - p.y, other.z - p.z}
}

func (p Point3D) String() string {
	return fmt.Sprintf("(%d,%d,%d)", p.x, p.y, p.z)
}
