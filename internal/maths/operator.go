package maths

import "math"

func PowInt(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}
