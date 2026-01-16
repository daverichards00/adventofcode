package maths

import (
	"golang.org/x/exp/constraints"
	"math"
)

func Abs[T constraints.Integer | constraints.Float](a T) T {
	if a < 0 {
		return -a
	}
	return a
}

func Factorial[T constraints.Integer](a T) T {
	for i := a - 1; i > 1; i, a = i-1, a*i {
	}
	return a
}

func PowInt(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}
