package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/internal/convert"
	"github.com/daverichards00/adventofcode/internal/file"
	"github.com/daverichards00/adventofcode/internal/maths"
	"github.com/daverichards00/adventofcode/internal/slice"
	"github.com/daverichards00/adventofcode/internal/spatial"
	"slices"
	"strings"
)

func main() {
	fmt.Println("Day 09")

	var tiles []spatial.Point2D

	input := file.Load("cmd/2025/day09/input.txt")
	for i := range input {
		parts := strings.Split(input[i], ",")
		tiles = append(tiles, spatial.NewPoint2D(convert.StrToInt(parts[0]), convert.StrToInt(parts[1])))
	}

	shape := make([]line, len(tiles))
	for i := 1; i < len(tiles); i++ {
		shape[i] = line{tiles[i-1], tiles[i]}
	}
	shape[0] = line{tiles[len(tiles)-1], tiles[0]}

	maxAreaA, maxAreaB := 0, 0
	for i := 0; i < len(tiles)-1; i++ {
		for j := i + 1; j < len(tiles); j++ {
			a := area(tiles[i], tiles[j])
			if a > maxAreaA {
				maxAreaA = a
			}
			if a > maxAreaB && areaWithinShape(tiles[i], tiles[j], shape) {
				maxAreaB = a
			}
		}
	}

	fmt.Println("Part A:")
	fmt.Printf("Maximum area between 2 points: %d\n\n", maxAreaA)

	fmt.Println("Part B:")
	fmt.Printf("Maximum area between 2 points within shape: %d\n\n", maxAreaB)
}

type line struct {
	a, b spatial.Point2D
}

func (l line) minX() int {
	return min(l.a.X(), l.b.X())
}
func (l line) minY() int {
	return min(l.a.Y(), l.b.Y())
}
func (l line) maxX() int {
	return max(l.a.X(), l.b.X())
}
func (l line) maxY() int {
	return max(l.a.Y(), l.b.Y())
}

func (l line) isVertical() bool {
	return l.a.X() == l.b.X()
}

func (l line) isHorizontal() bool {
	return l.a.Y() == l.b.Y()
}

func (l line) intersects(other line) bool {
	if l.isVertical() == other.isVertical() {
		// Parallel
		return false
	}
	if l.isVertical() {
		return other.a.Y() > l.minY() && other.a.Y() < l.maxY() && other.minX() < l.a.X() && other.maxX() > l.a.X()
	}
	return other.a.X() > l.minX() && other.a.X() < l.maxX() && other.minY() < l.a.Y() && other.maxY() > l.a.Y()
}

func (l line) overlaps(other line) bool {
	if l.isVertical() && other.isVertical() {
		return other.a.X() == l.a.X() && other.minY() <= l.maxY() && other.maxY() >= l.minY()
	}
	if l.isHorizontal() && other.isHorizontal() {
		return other.a.Y() == l.a.Y() && other.minX() <= l.maxX() && other.maxX() >= l.minX()
	}
	return false
}

func (l line) contains(p spatial.Point2D) bool {
	// Assume lines are only vertical or horizontal
	return p.X() >= l.minX() && p.X() <= l.maxX() && p.Y() >= l.minY() && p.Y() <= l.maxY()
}

func area(a, b spatial.Point2D) int {
	return (maths.Abs(a.X()-b.X()) + 1) * (maths.Abs(a.Y()-b.Y()) + 1)
}

func areaWithinShape(a, b spatial.Point2D, shape []line) bool {
	// Work out if this area (a->b) is entirely within the shape provided as a slice of boundary lines.
	// If each line of the area is within the shape, then the whole area must be within the shape.

	// Check boundary lines
	bLines := []line{
		{a, spatial.NewPoint2D(a.X(), b.Y())},
		{spatial.NewPoint2D(b.X(), a.Y()), b},
		{a, spatial.NewPoint2D(b.X(), a.Y())},
		{spatial.NewPoint2D(a.X(), b.Y()), b},
	}
	bLines = slice.Unique(slice.Filter(bLines, func(l line) bool { return l.a != l.b }))

	for _, bl := range bLines {
		if !lineWithinShape(bl, shape) {
			return false
		}
	}
	return true
}

func lineWithinShape(l line, shape []line) bool {
	var overlapping []line
	for _, sl := range shape {
		// If intersecting, line is partially out of bounds
		if l.intersects(sl) {
			return false
		}
		// If overlapping, make sure we don't test a point on this line
		if l.overlaps(sl) {
			overlapping = append(overlapping, sl)
		}
	}
	// Check a point is within the shape
	p := spatial.NewPoint2D(l.minX(), l.minY())
	v := p.To(spatial.NewPoint2D(l.maxX(), l.maxY())).Unit()
	for slices.ContainsFunc(overlapping, func(ol line) bool {
		return ol.contains(p)
	}) {
		p = p.Add(v)
		if !l.contains(p) {
			// All points on a line
			return true
		}
	}

	return pointWithinShape(p, shape)
}

func pointWithinShape(p spatial.Point2D, shape []line) bool {
	var linesToLeft []line
	for i, sl := range shape {
		if sl.contains(p) {
			// On a line (counts as within bounds)
			return true
		}
		if sl.isVertical() {
			if sl.a.X() < p.X() && sl.contains(spatial.NewPoint2D(sl.a.X(), p.Y())) {
				// Vertical line to the left
				linesToLeft = append(linesToLeft, sl)
			}
			continue
		}
		// Horizontal
		if sl.a.Y() != p.Y() {
			// Not on the same Y coord
			continue
		}
		if sl.minX() > p.X() {
			// To the right
			continue
		}
		// If the 2 lines connected to this horizontal line are going in the same direction, include
		// (We want to count `s` and `z` lines as a boundary, but not `n` and `u` shapes)
		prev, next := sliceGetWrapped(shape, i-1), sliceGetWrapped(shape, i+1)
		if (prev.minY() < p.Y()) != (next.minY() < p.Y()) {
			linesToLeft = append(linesToLeft, sl)
		}
	}
	// Odd number of lines => within
	// Even number of lines => outside
	return len(linesToLeft)%2 > 0
}

func sliceGetWrapped[T comparable](s []T, idx int) T {
	l := len(s)
	for idx < 0 {
		idx += l
	}
	for idx >= l {
		idx -= l
	}
	return s[idx]
}
