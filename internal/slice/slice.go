package slice

import (
	"golang.org/x/exp/constraints"
	"slices"
)

func Filter[T any](s []T, f func(T) bool) []T {
	var r []T
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}

func IndexAllFunc[T any](s []T, f func(T) bool) []int {
	var r []int
	for i, v := range s {
		if f(v) {
			r = append(r, i)
		}
	}
	return r
}

func Intersect[T comparable](a, b []T) []T {
	var r []T
	for _, v := range a {
		if slices.Contains(b, v) {
			r = append(r, v)
		}
	}
	return r
}

func Reduce[T any](s []T, f func(T, T) T) T {
	var r T
	if len(s) == 0 {
		return r
	}
	r = s[0]
	for _, v := range s[1:] {
		r = f(r, v)
	}
	return r
}

func Split[T comparable](s []T, v T) [][]T {
	var r [][]T
	for {
		j := slices.Index(s, v)
		if j < 0 {
			break
		}
		r = append(r, s[:j])
		s = s[j+1:]
	}
	r = append(r, s)
	return r
}

func Sum[T constraints.Integer | constraints.Float](s []T) T {
	var r T
	for _, v := range s {
		r += v
	}
	return r
}

func Unique[T comparable](s []T) []T {
	var r []T
	for _, v := range s {
		if !slices.Contains(r, v) {
			r = append(r, v)
		}
	}
	return r
}
