package spatial

import (
	"fmt"
	"github.com/daverichards00/adventofcode/internal/maths"
)

type Vector3D struct {
	x, y, z int
}

func NewVector3D(x, y, z int) Vector3D {
	return Vector3D{x, y, z}
}

func (v Vector3D) X() int {
	return v.x
}

func (v Vector3D) Y() int {
	return v.y
}

func (v Vector3D) Z() int {
	return v.z
}

func (v Vector3D) Manhattan() int {
	return maths.Abs(v.x) + maths.Abs(v.y) + maths.Abs(v.z)
}

func (v Vector3D) String() string {
	return fmt.Sprintf("(%d,%d,%d)", v.x, v.y, v.z)
}
