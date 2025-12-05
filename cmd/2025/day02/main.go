package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/internal/convert"
	"github.com/daverichards00/adventofcode/internal/file"
	"github.com/daverichards00/adventofcode/internal/maths"
	"github.com/daverichards00/adventofcode/internal/slice"
	"strings"
)

func main() {
	fmt.Println("Day 02")

	input := file.Load("cmd/2025/day02/input.txt")
	ranges := parseRanges(input[0])

	partA, partB := 0, 0
	for _, r := range ranges {
		partA += slice.Sum(r.invalidsN(2))
		partB += slice.Sum(r.invalids())
	}

	fmt.Println("Part A:")
	fmt.Printf("The sum of invalid IDs: %d\n\n", partA)

	fmt.Println("Part B:")
	fmt.Printf("The sum of invalid IDs: %d\n\n", partB)
}

func parseRanges(s string) []rng {
	var ranges []rng
	for _, line := range strings.Split(s, ",") {
		parts := strings.Split(line, "-")
		if len(parts) != 2 {
			panic("bad input")
		}
		ranges = append(ranges, rng{
			start: convert.StrToInt(parts[0]),
			end:   convert.StrToInt(parts[1]),
		})
	}
	return ranges
}

type rng struct {
	start, end int
}

func (r rng) invalids() []int {
	var invalids []int

	maxDC := digitCount(r.end)
	for n := 2; n <= maxDC; n++ {
		invalids = append(invalids, r.invalidsN(n)...)
	}

	return slice.Unique(invalids)
}

func (r rng) invalidsN(n int) []int {
	// n: number of repeated parts

	s := r.start
	sdc := digitCount(s)
	if sdc%n != 0 {
		// round up to next number with digit count multiple of n
		s = maths.PowInt(10, (((sdc/n)+1)*n)-1)
		sdc = digitCount(s)
	}
	sFrag := s / maths.PowInt(10, sdc-(sdc/n))
	sFull := repeatDigitsN(sFrag, n)

	if sFull < s {
		sFrag++
	}

	e := r.end
	edc := digitCount(e)
	if edc%n != 0 {
		// round down to next number with even digit count
		e = maths.PowInt(10, (edc/n)*n) - 1
		edc = digitCount(e)
	}
	eFrag := e / maths.PowInt(10, edc-(edc/n))
	eFull := repeatDigitsN(eFrag, n)

	if eFull > e {
		eFrag--
	}

	var invalid []int
	for frag := sFrag; frag <= eFrag; frag++ {
		invalid = append(invalid, repeatDigitsN(frag, n))
	}

	return invalid
}

func digitCount(i int) int {
	if i < 0 {
		return 0
	}
	for p := 1; p < 20; p++ {
		if i < maths.PowInt(10, p) {
			return p
		}
	}
	panic("number too big")
}

func repeatDigitsN(frag, n int) int {
	dc := digitCount(frag)
	r := frag
	for i := 1; i < n; i++ {
		r = (r * maths.PowInt(10, dc)) + frag
	}
	return r
}
