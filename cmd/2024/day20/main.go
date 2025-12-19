package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/internal/file"
	"github.com/daverichards00/adventofcode/internal/pathfinder"
	"github.com/daverichards00/adventofcode/internal/spatial"
	"slices"
)

func main() {
	fmt.Println("Day 20")

	s := spatial.NewSpace2D[bool]()
	var start spatial.Point2D
	var end spatial.Point2D

	input := file.Load("cmd/2024/day20/input.txt")
	for y, line := range input {
		for x, r := range line {
			p := spatial.NewPoint2D(x, y)
			switch r {
			case '#':
				s.Set(p, true)
			case 'S':
				start = p
			case 'E':
				end = p
			}
		}
	}

	bestNonCheatPath := findNonCheat(s, start, end)

	toSave := 100

	fmt.Println("Part A:")
	partACheats := findAllCheats(bestNonCheatPath, 2, toSave)
	fmt.Printf("Number of cheats that save at least %d picoseconds: %d\n\n", toSave, partACheats)

	fmt.Println("Part B:")
	partBCheats := findAllCheats(bestNonCheatPath, 20, toSave)
	fmt.Printf("Number of cheats that save at least %d picoseconds: %d\n\n", toSave, partBCheats)
}

var directions = []spatial.Vector2D{spatial.North2D, spatial.South2D, spatial.East2D, spatial.West2D}

func findNonCheat(s *spatial.Space2D[bool], start, end spatial.Point2D) []spatial.Point2D {
	comp, err := pathfinder.Find(
		[][]spatial.Point2D{{start}},
		pathfinder.Fn[[]spatial.Point2D]{
			Next: func(p []spatial.Point2D) [][]spatial.Point2D {
				var next [][]spatial.Point2D
				for _, d := range directions {
					// Next step
					n := p[len(p)-1].Add(d)

					if s.Exists(n) || slices.Contains(p, n) {
						continue
					}

					// Add next path
					next = append(next, append(slices.Clone(p), n))
				}
				return next
			},
			Complete: func(p []spatial.Point2D) bool {
				return p[len(p)-1] == end
			},
			Less: func(a, b []spatial.Point2D) bool {
				return len(a) < len(b)
			},
		},
	)
	if err != nil {
		panic(err)
	}
	return comp[0]
}

func findAllCheats(path []spatial.Point2D, maxCheatDist, minSave int) int {
	cheatCount := 0
	for i := 0; i < len(path)-minSave; i++ {
		// Find all points on the path within the maxCheat distance
		for j := i + minSave; j < len(path); j++ {
			cheatDist := path[i].To(path[j]).Manhattan()
			if cheatDist > maxCheatDist {
				// Can't cheat this far
				continue
			}
			if i+cheatDist <= j-minSave {
				// This cheat would save enough steps
				cheatCount++
			}
		}
	}
	return cheatCount
}
