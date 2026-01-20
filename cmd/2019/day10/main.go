package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/internal/file"
	"github.com/daverichards00/adventofcode/internal/maths"
	"github.com/daverichards00/adventofcode/internal/spatial"
	"slices"
	"sort"
)

func main() {
	fmt.Println("Day 10")

	var asteroids []spatial.Point2D

	input := file.Load("cmd/2019/day10/input.txt")
	for y, line := range input {
		for x, r := range line {
			if r == '#' {
				asteroids = append(asteroids, spatial.NewPoint2D(x, y))
			}
		}
	}

	fmt.Println("Part A:")
	bestLocation, bestDetectable := 0, detectable(asteroids, asteroids[0])
	for i := 1; i < len(asteroids); i++ {
		if d := detectable(asteroids, asteroids[i]); len(d) > len(bestDetectable) {
			bestLocation, bestDetectable = i, d
		}
	}
	fmt.Printf("Number of asteroids detectable for best location: %d\n\n", len(bestDetectable))

	fmt.Println("Part B:")
	destroyed := sortLaser(bestDetectable, asteroids[bestLocation])
	fmt.Printf("200th asteroid to be destroyed: %d\n\n", destroyed[199].X()*100+destroyed[199].Y())
}

func detectable(asteroids []spatial.Point2D, location spatial.Point2D) []spatial.Point2D {
	var vectors []spatial.Vector2D
	for _, a := range asteroids {
		if location == a {
			continue
		}
		v := location.To(a).Min()
		if !slices.Contains(vectors, v) {
			vectors = append(vectors, v)
		}
	}
	var d []spatial.Point2D
	for _, v := range vectors {
		n := location.Add(v)
		for {
			if slices.Contains(asteroids, n) {
				d = append(d, n)
				break
			}
			n = n.Add(v)
		}
	}
	return d
}

func sortLaser(asteroids []spatial.Point2D, location spatial.Point2D) []spatial.Point2D {
	vectors := make([]spatial.Vector2D, len(asteroids))
	for i := range asteroids {
		vectors[i] = location.To(asteroids[i])
	}
	sort.Slice(vectors, func(i, j int) bool {
		a, b := vectors[i], vectors[j]
		if (a.X() >= 0) != (b.X() >= 0) {
			return a.X() > b.X()
		}
		if (a.Y() >= 0) != (b.Y() >= 0) {
			if a.X() >= 0 {
				return a.Y() < b.Y()
			}
			return a.Y() > b.Y()
		}
		// Within same quadrant, compare angles
		var aO, aA, bO, bA float64
		if (a.X() >= 0) == (a.Y() >= 0) {
			// Opposite => Y, Adjacent = X
			aO, aA, bO, bA = maths.Abs(float64(a.Y())), maths.Abs(float64(a.X())), maths.Abs(float64(b.Y())), maths.Abs(float64(b.X()))
		} else {
			// Opposite => X, Adjacent = Y
			aO, aA, bO, bA = maths.Abs(float64(a.X())), maths.Abs(float64(a.Y())), maths.Abs(float64(b.X())), maths.Abs(float64(b.Y()))
		}
		return aO/aA < bO/bA
	})

	o := make([]spatial.Point2D, len(vectors))
	for i := range vectors {
		o[i] = location.Add(vectors[i])
	}
	return o
}
