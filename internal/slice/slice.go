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
