package slice

import (
	"golang.org/x/exp/constraints"
	"slices"
)

func Sum[T constraints.Integer | constraints.Float](s []T) T {
	var sm T
	for _, v := range s {
		sm += v
	}
	return sm
}

func Unique[T comparable](s []T) []T {
	var u []T
	for _, v := range s {
		if !slices.Contains(u, v) {
			u = append(u, v)
		}
	}
	return u
}
