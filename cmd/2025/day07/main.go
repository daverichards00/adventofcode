package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/internal/file"
	"github.com/daverichards00/adventofcode/internal/slice"
	"github.com/daverichards00/adventofcode/internal/spatial"
	"maps"
	"slices"
)

func main() {
	fmt.Println("Day 07")

	space := spatial.NewSpace2D[bool]()
	var start spatial.Point2D

	input := file.Load("cmd/2025/day07/input.txt")
	for y, line := range input {
		for x, c := range line {
			if c == 'S' {
				// start
				start = spatial.NewPoint2D(x, y)
			}
			if c == '^' {
				// splitter
				space.Set(spatial.NewPoint2D(x, y), true)
			}
		}
	}

	fmt.Println("Part A:")
	splitFreq := 0
	beams := []spatial.Point2D{start}
	for range input {
		var next []spatial.Point2D
		for _, b := range beams {
			n := b.Add(spatial.South2D)
			if space.GetOrDefault(n, false) {
				next = append(next, n.Add(spatial.West2D), n.Add(spatial.East2D))
				splitFreq++
				continue
			}
			next = append(next, n)
		}
		beams = slice.Unique(next)
	}
	fmt.Printf("Number of times beam split: %d\n\n", splitFreq)

	fmt.Println("Part B:")
	beamTimes := map[spatial.Point2D]int{start: 1}
	for range input {
		next := map[spatial.Point2D]int{}
		for b, t := range beamTimes {
			n := b.Add(spatial.South2D)
			if space.GetOrDefault(n, false) {
				w, e := n.Add(spatial.West2D), n.Add(spatial.East2D)
				next[w] = next[w] + t
				next[e] = next[e] + t
				continue
			}
			next[n] = next[n] + t
		}
		beamTimes = next
	}
	partB := slice.Sum(slices.Collect(maps.Values(beamTimes)))
	fmt.Printf("Number of timelines: %d\n\n", partB)
}
