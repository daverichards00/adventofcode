package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/internal/convert"
	"github.com/daverichards00/adventofcode/internal/file"
	"github.com/daverichards00/adventofcode/internal/spatial"
	"slices"
	"strings"
)

func main() {
	fmt.Println("Day 08")

	var jboxes []spatial.Point3D

	input := file.Load("cmd/2025/day08/input.txt")
	for _, line := range input {
		parts := strings.Split(line, ",")
		jboxes = append(jboxes, spatial.NewPoint3D(convert.StrToInt(parts[0]), convert.StrToInt(parts[1]), convert.StrToInt(parts[2])))
	}

	var conns []conn
	for i := 0; i < len(jboxes)-1; i++ {
		for j := i + 1; j < len(jboxes); j++ {
			conns = append(conns, newConn(jboxes[i], jboxes[j]))
		}
	}

	slices.SortFunc(conns, func(a, b conn) int {
		return a.sqDist - b.sqDist
	})

	fmt.Println("Part A:")
	shortestConns := conns[0:1000]
	circuits, _ := buildCircuits(jboxes, shortestConns)

	slices.SortFunc(circuits, func(a, b []spatial.Point3D) int {
		return len(b) - len(a)
	})

	fmt.Printf("Product of the size of the 3 largest circuits: %d\n\n", len(circuits[0])*len(circuits[1])*len(circuits[2]))

	fmt.Println("Part B:")
	_, lastConn := buildCircuits(jboxes, conns)
	fmt.Printf("Product of X coordinates of last connection to make 1 circuit: %d\n\n", lastConn.a.X()*lastConn.b.X())
}

type conn struct {
	a, b   spatial.Point3D
	sqDist int
}

func newConn(a, b spatial.Point3D) conn {
	return conn{a, b, (a.X()-b.X())*(a.X()-b.X()) + (a.Y()-b.Y())*(a.Y()-b.Y()) + (a.Z()-b.Z())*(a.Z()-b.Z())}
}

func buildCircuits(jboxes []spatial.Point3D, conns []conn) ([][]spatial.Point3D, conn) {
	var circuits [][]spatial.Point3D
	var lastConnIdx int
	for i, c := range conns {
		lastConnIdx = i
		aIdx := slices.IndexFunc(circuits, func(circuit []spatial.Point3D) bool { return slices.Contains(circuit, c.a) })
		bIdx := slices.IndexFunc(circuits, func(circuit []spatial.Point3D) bool { return slices.Contains(circuit, c.b) })
		if aIdx < 0 && bIdx < 0 {
			// Neither jboxes in a circuit yet
			circuits = append(circuits, []spatial.Point3D{c.a, c.b})
			continue
		}
		if aIdx < 0 {
			// B found, add A in
			circuits[bIdx] = append(circuits[bIdx], c.a)
			continue
		}
		if bIdx < 0 {
			// A found, add B in
			circuits[aIdx] = append(circuits[aIdx], c.b)
			continue
		}
		// Both found
		if aIdx == bIdx {
			// Already in same circuit, do nothing
			continue
		}
		// Need to merge
		// Add B to A
		circuits[aIdx] = append(circuits[aIdx], circuits[bIdx]...)
		// Remove B
		circuits = slices.Delete(circuits, bIdx, bIdx+1)
		// One circuit?
		if len(circuits) == 1 && len(circuits[0]) == len(jboxes) {
			break
		}
	}
	// Add remaining jboxes as circuits of 1
	for _, jbox := range jboxes {
		if slices.ContainsFunc(circuits, func(circuit []spatial.Point3D) bool { return slices.Contains(circuit, jbox) }) {
			continue
		}
		circuits = append(circuits, []spatial.Point3D{jbox})
	}
	return circuits, conns[lastConnIdx]
}
