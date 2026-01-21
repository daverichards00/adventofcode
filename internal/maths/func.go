package maths

import (
	"golang.org/x/exp/constraints"
	"math"
	"slices"
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

func LowestCommonMultiple[T constraints.Integer](a T, b ...T) T {
	f := PrimeFactors(a)
	for i := range b {
		ff := slices.Clone(f)
		for _, bf := range PrimeFactors(b[i]) {
			if idx := slices.Index(ff, bf); idx >= 0 {
				ff = slices.Delete(ff, idx, idx+1)
				continue
			}
			f = append(f, bf)
		}
	}
	lcm := f[0]
	for i := 1; i < len(f); i++ {
		lcm *= f[i]
	}
	return lcm
}

func PowInt(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}

func PrimeFactors[T constraints.Integer](a T) []T {
	if a < 0 {
		panic("PrimeFactors called with negative value")
	}
	if a < 4 {
		return []T{a}
	}
	var f []T
PrimeFactors:
	for {
		// Treat 2 separately, so we can step the following loop by 2 instead of 1
		if a == 2 {
			return append(f, 2)
		}
		if a%2 == 0 {
			f, a = append(f, 2), a/2
			continue PrimeFactors
		}
		for i := T(3); (i * i) <= a; i += 2 {
			if a%i == 0 {
				f, a = append(f, i), a/i
				continue PrimeFactors
			}
		}
		return append(f, a)
	}
}
