package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/internal/file"
	"github.com/daverichards00/adventofcode/internal/slice"
	"github.com/daverichards00/adventofcode/internal/spatial"
)

func main() {
	fmt.Println("Day 04")

	space := spatial.NewSpace2D[bool]()

	input := file.Load("cmd/2025/day04/input.txt")
	for y, line := range input {
		for x, r := range line {
			if r == '@' {
				space.Set(spatial.NewPoint2D(x, y), true)
			}
		}
	}

	fmt.Println("Part A:")
	accessible := getAccessible(space)
	fmt.Printf("Number of accessible rolls: %d\n\n", len(accessible))

	fmt.Println("Part B:")
	removed := 0
	for {
		acc := getAccessible(space)
		if len(acc) == 0 {
			break
		}
		for _, p := range acc {
			space.Unset(p)
			removed++
		}
	}
	fmt.Printf("Total removed: %d\n\n", removed)

}

func getAccessible(space *spatial.Space2D[bool]) []spatial.Point2D {
	return slice.Filter(space.Find(true), func(p spatial.Point2D) bool {
		return len(space.GetAllAdjacentAndDiagonal(p)) < 4
	})
}
