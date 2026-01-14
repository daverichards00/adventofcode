package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/internal/convert"
	"github.com/daverichards00/adventofcode/internal/file"
	"github.com/daverichards00/adventofcode/internal/spatial"
	"sort"
	"strings"
)

func main() {
	fmt.Println("Day 03")

	input := file.Load("cmd/2019/day03/input.txt")
	wireA, wireB := strToWire(input[0]), strToWire(input[1])

	var origin = spatial.NewPoint2D(0, 0)

	var intersects []intersect
	stepsA := 0
	for _, lA := range wireA {
		stepsB := 0
		for _, lB := range wireB {
			if i, ok := findIntersect(lA, lB); ok {
				intersects = append(intersects, intersect{
					p:      i,
					m:      origin.To(i).Manhattan(),
					stepsA: stepsA + lA.p.To(i).Manhattan(),
					stepsB: stepsB + lB.p.To(i).Manhattan(),
				})
			}
			stepsB += lB.v.Manhattan()
		}
		stepsA += lA.v.Manhattan()
	}

	fmt.Println("Part A:")
	sort.Slice(intersects, func(i, j int) bool { return intersects[i].m < intersects[j].m })
	fmt.Printf("Distance of closest intersection to port: %d\n\n", intersects[0].m)

	fmt.Println("Part B:")
	sort.Slice(intersects, func(i, j int) bool {
		return (intersects[i].stepsA + intersects[i].stepsB) < (intersects[j].stepsA + intersects[j].stepsB)
	})
	fmt.Printf("Fewest combined steps to intersection: %d\n\n", intersects[0].stepsA+intersects[0].stepsB)
}

type line struct {
	p spatial.Point2D
	v spatial.Vector2D
}

func strToWire(s string) []line {
	var w []line
	p := spatial.NewPoint2D(0, 0)
	for _, d := range strings.Split(s, ",") {
		v := strToVector(d)
		w = append(w, line{
			p: p,
			v: v,
		})
		p = p.Add(v)
	}
	return w
}

func strToVector(s string) spatial.Vector2D {
	switch s[0] {
	case 'U':
		return spatial.NewVector2D(0, convert.StrToInt(s[1:]))
	case 'D':
		return spatial.NewVector2D(0, -convert.StrToInt(s[1:]))
	case 'L':
		return spatial.NewVector2D(-convert.StrToInt(s[1:]), 0)
	case 'R':
		return spatial.NewVector2D(convert.StrToInt(s[1:]), 0)
	default:
		panic("Unknown direction")
	}
}

type intersect struct {
	p      spatial.Point2D
	m      int
	stepsA int
	stepsB int
}

func findIntersect(a, b line) (spatial.Point2D, bool) {
	if (a.v.X() == 0) == (b.v.X() == 0) {
		// Parallel
		return spatial.Point2D{}, false
	}
	if a.v.X() == 0 {
		a, b = b, a
	}
	// a: horizontal, b: vertical
	if (a.p.X() < b.p.X()) != (a.p.Add(a.v).X() < b.p.X()) && (b.p.Y() < a.p.Y()) != (b.p.Add(b.v).Y() < a.p.Y()) {
		return spatial.NewPoint2D(b.p.X(), a.p.Y()), true
	}
	return spatial.Point2D{}, false
}
