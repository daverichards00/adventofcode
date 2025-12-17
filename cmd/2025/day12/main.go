package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/internal/convert"
	"github.com/daverichards00/adventofcode/internal/file"
	"github.com/daverichards00/adventofcode/internal/slice"
	"strings"
)

var (
	// All presents bound within 3x3
	presentBound = 3
)

func main() {
	fmt.Println("Day 12")

	var presents []present
	var regions []region

	input := file.Load("cmd/2025/day12/input.txt")
	inputParts := slice.Split(input, "")

	for _, ip := range inputParts[:len(inputParts)-1] {
		s, v := make([][]bool, presentBound), 0
		for i, line := range ip[1:] {
			s[i] = make([]bool, presentBound)
			for j, r := range line {
				if r == '#' {
					s[i][j] = true
					v++
				}
			}
		}
		presents = append(presents, present{shape: s, volume: v})
	}

	for _, r := range inputParts[len(inputParts)-1] {
		p := strings.Split(r, ":")
		wl := strings.Split(p[0], "x")
		q := strings.Split(strings.TrimSpace(p[1]), " ")

		rg := region{w: convert.StrToInt(wl[0]), l: convert.StrToInt(wl[1]), quantities: make([]int, len(q))}
		for i := range q {
			rg.quantities[i] = convert.StrToInt(q[i])
		}

		regions = append(regions, rg)
	}

	fmt.Println("Part A:")
	partA := 0
	for _, r := range regions {
		if canFit(r, presents) {
			partA++
		}
	}
	fmt.Printf("Number of regions that can fit all presents: %d\n\n", partA)
}

type present struct {
	shape  [][]bool
	volume int
}

type region struct {
	w, l       int
	quantities []int
}

func canFit(r region, presents []present) bool {
	if (r.w/presentBound)*(r.l/presentBound) >= slice.Sum(r.quantities) {
		// Can definitely fit
		return true
	}
	minVol := 0
	for i, q := range r.quantities {
		minVol += q * presents[i].volume
	}
	if r.w*r.l < minVol {
		// Definitely not possible
		return false
	}
	// Maybe?
	panic("needs implementing")
}
