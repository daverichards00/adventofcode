package spatial

import (
	"fmt"
	"github.com/daverichards00/adventofcode/internal/maths"
)

// North set as -1 and South set as +1 in the y-axis to accommodate puzzle input generally being read from top to bottom
// (the line number (y-axis) of the input increasing as it is read downwards)
var (
	North2D = Vector2D{0, -1}
	South2D = Vector2D{0, 1}
	West2D  = Vector2D{-1, 0}
	East2D  = Vector2D{1, 0}

	NorthWest2D = Vector2D{-1, -1}
	NorthEast2D = Vector2D{1, -1}
	SouthWest2D = Vector2D{-1, 1}
	SouthEast2D = Vector2D{1, 1}
)

type Vector2D struct {
	x, y int
}

func NewVector2D(x, y int) Vector2D {
	return Vector2D{x, y}
}

func (v Vector2D) X() int {
	return v.x
}

func (v Vector2D) Y() int {
	return v.y
}

func (v Vector2D) Unit() Vector2D {
	// Only true unit vector for the 4 compass directions
	u := Vector2D{v.x, v.y}
	if u.x != 0 {
		u.x /= maths.Abs(u.x)
	}
	if u.y != 0 {
		u.y /= maths.Abs(u.y)
	}
	return u
}

func (v Vector2D) String() string {
	return fmt.Sprintf("(%d,%d)", v.x, v.y)
}
