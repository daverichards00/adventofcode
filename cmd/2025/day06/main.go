package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/internal/convert"
	"github.com/daverichards00/adventofcode/internal/file"
	"github.com/daverichards00/adventofcode/internal/slice"
	"regexp"
	"strings"
)

var colSeparator = regexp.MustCompile(` +`)

func main() {
	fmt.Println("Day 06")

	input := file.Load("cmd/2025/day06/input.txt")
	eqs := parseInput(input)

	fmt.Println("Part A:")
	var partA []int
	for _, e := range eqs {
		partA = append(partA, e.calcA())
	}
	fmt.Printf("The sum of calculations: %d\n\n", slice.Sum(partA))

	fmt.Println("Part B:")
	var partB []int
	for _, e := range eqs {
		partB = append(partB, e.calcB())
	}
	fmt.Printf("The sum of calculations: %d\n\n", slice.Sum(partB))

}

func parseInput(input []string) []eq {
	// cols
	sep := columnSeparatorIndices(input)
	longestLen := len(longest(input))
	// vals
	eqs := make([]eq, len(sep)+1)
	for i := range eqs {
		eqs[i].vals = make([]string, len(input)-1)
	}
	for i, line := range input[0 : len(input)-1] {
		// Pad to longest
		line = line + strings.Repeat(" ", longestLen-len(line))
		prev := -1
		for j, idx := range sep {
			eqs[j].vals[i] = line[prev+1 : idx]
			prev = idx
		}
		eqs[len(sep)].vals[i] = line[prev+1:]
	}
	// ops
	parts := colSeparator.Split(strings.TrimSpace(input[len(input)-1]), -1)
	for i := range parts {
		eqs[i].op = parts[i]
	}
	return eqs
}

func columnSeparatorIndices(input []string) []int {
	indices := StrFindAllIndices(input[0], ' ')
	for i := range input[1:] {
		indices = slice.Intersect(indices, StrFindAllIndices(input[i], ' '))
	}
	return indices
}

func StrFindAllIndices(s string, c rune) []int {
	var r []int
	for i := range s {
		if rune(s[i]) == c {
			r = append(r, i)
		}
	}
	return r
}

type eq struct {
	vals []string
	op   string
}

func (e eq) valsA() []int {
	v := make([]int, len(e.vals))
	for i := range e.vals {
		v[i] = convert.StrToInt(strings.TrimSpace(e.vals[i]))
	}
	return v
}

func (e eq) valsB() []int {
	v := make([]int, len(e.vals[0]))
	for i := range e.vals[0] {
		s := ""
		for j := range e.vals {
			s += string(e.vals[j][i])
		}
		v[i] = convert.StrToInt(strings.TrimSpace(s))
	}
	return v
}

func (e eq) calcA() int {
	return e.calc(e.valsA())
}

func (e eq) calcB() int {
	return e.calc(e.valsB())
}

func (e eq) calc(vals []int) int {
	switch e.op {
	case "+":
		return slice.Reduce(vals, func(a int, b int) int {
			return a + b
		})
	case "*":
		return slice.Reduce(vals, func(a int, b int) int {
			return a * b
		})
	default:
		panic("invalid op")
	}
}

func longest(s []string) string {
	return slice.Reduce(s, func(a string, b string) string {
		if len(b) > len(a) {
			return b
		}
		return a
	})
}
