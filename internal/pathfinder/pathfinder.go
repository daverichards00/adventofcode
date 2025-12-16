package pathfinder

import (
	"errors"
	"sort"
)

type Fn[T any] struct {
	// Next takes a path and returns the possible new path(s) with the next step(s) added on.
	Next func(T) []T
	// Complete determines whether a path is complete.
	Complete func(T) bool
	// Less is used to sort paths in order of how optimal they are. True should be returned if path a is better than path b.
	Less func(a, b T) bool
}

func Find[T any](initial []T, fn Fn[T]) ([]T, error) {
	paths := initial
	for len(paths) > 0 {
		// pop shortest
		path := paths[0]
		if len(paths) == 1 {
			paths = nil
		} else {
			paths = paths[1:]
		}

		if fn.Complete(path) {
			// if complete, return all complete with equal score
			comp := []T{path}
			for _, p := range paths {
				if fn.Less(path, p) {
					break
				}
				if !fn.Complete(p) {
					continue
				}
				comp = append(comp, p)
			}
			return comp, nil
		}

		// add possible next paths
		next := fn.Next(path)

		// if new paths found, add and sort
		if len(next) > 0 {
			paths = append(paths, next...)
			sort.Slice(paths, func(i, j int) bool {
				return fn.Less(paths[i], paths[j])
			})
		}
	}
	
	return nil, errors.New("no paths found")
}
