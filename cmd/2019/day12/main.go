package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/internal/convert"
	"github.com/daverichards00/adventofcode/internal/file"
	"github.com/daverichards00/adventofcode/internal/maths"
	"github.com/daverichards00/adventofcode/internal/slice"
	"github.com/daverichards00/adventofcode/internal/spatial"
	"golang.org/x/exp/constraints"
	"regexp"
	"slices"
)

var inputRegExp = regexp.MustCompile(`^<x=(-?[0-9]+), y=(-?[0-9]+), z=(-?[0-9]+)>$`)

func main() {
	fmt.Println("Day 12")

	var moons []object

	input := file.Load("cmd/2019/day12/input.txt")
	for _, line := range input {
		m := inputRegExp.FindStringSubmatch(line)
		moons = append(moons, object{
			p: spatial.NewPoint3D(convert.StrToInt(m[1]), convert.StrToInt(m[2]), convert.StrToInt(m[3])),
			v: spatial.NewVector3D(0, 0, 0),
		})
	}

	fmt.Println("Part A:")
	moonsA := slices.Clone(moons)
	for i := 0; i < 1000; i++ {
		step(moonsA)
	}
	partA := 0
	for _, moon := range moonsA {
		partA += moon.energy()
	}
	fmt.Printf("Total energy in the system: %d\n\n", partA)

	fmt.Println("Part B:")
	stepsX := stepsBeforeLoop(slice.Transform(moons, func(m object) int { return m.p.X() }), slice.Transform(moons, func(m object) int { return m.v.X() }))
	stepsY := stepsBeforeLoop(slice.Transform(moons, func(m object) int { return m.p.Y() }), slice.Transform(moons, func(m object) int { return m.v.Y() }))
	stepsZ := stepsBeforeLoop(slice.Transform(moons, func(m object) int { return m.p.Z() }), slice.Transform(moons, func(m object) int { return m.v.Z() }))
	fmt.Printf("Minimum steps to reach a previous state: %d\n\n", maths.LowestCommonMultiple(stepsX, stepsY, stepsZ))
}

type object struct {
	p spatial.Point3D
	v spatial.Vector3D
}

func (o object) energy() int {
	return spatial.NewPoint3D(0, 0, 0).To(o.p).Manhattan() * o.v.Manhattan()
}

func (o object) equal(other object) bool {
	return o.p == other.p && o.v == other.v
}

func step(moons []object) {
	// Gravity
	for i := 0; i < len(moons)-1; i++ {
		for j := i + 1; j < len(moons); j++ {
			d := moons[i].p.To(moons[j].p)
			moons[i].v = spatial.NewVector3D(moons[i].v.X()+norm(d.X()), moons[i].v.Y()+norm(d.Y()), moons[i].v.Z()+norm(d.Z()))
			moons[j].v = spatial.NewVector3D(moons[j].v.X()-norm(d.X()), moons[j].v.Y()-norm(d.Y()), moons[j].v.Z()-norm(d.Z()))
		}
	}
	// Move
	for i := range moons {
		moons[i].p = moons[i].p.Add(moons[i].v)
	}
}

func norm[T constraints.Integer](i T) T {
	if i == 0 {
		return i
	}
	return i / maths.Abs(i)
}

func stepsBeforeLoop(positions, vectors []int) int {
	initP, initV := slices.Clone(positions), slices.Clone(vectors)
	steps := 0
	for {
		// Gravity
		for i := 0; i < len(positions)-1; i++ {
			for j := i + 1; j < len(positions); j++ {
				d := positions[i] - positions[j]
				if d > 0 {
					vectors[i]--
					vectors[j]++
				} else if d < 0 {
					vectors[i]++
					vectors[j]--
				}
			}
		}
		// Move
		for i := range positions {
			positions[i] += vectors[i]
		}
		steps++
		if slices.Equal(initP, positions) && slices.Equal(initV, vectors) {
			return steps
		}
	}
}
