package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/internal/convert"
	"github.com/daverichards00/adventofcode/internal/file"
	"slices"
	"strings"
)

func main() {
	fmt.Println("Day 05")

	var ingRanges []ingRange
	var ingIDs []int

	input := file.Load("cmd/2025/day05/input.txt")

	splitIdx := slices.Index(input, "")
	for _, r := range input[0:splitIdx] {
		ingRanges = append(ingRanges, newIngRange(r))
	}
	for _, id := range input[splitIdx+1:] {
		ingIDs = append(ingIDs, convert.StrToInt(id))
	}

	ingRanges = normaliseIngRanges(ingRanges)

	fmt.Println("Part A:")
	partA := 0
	for _, i := range ingIDs {
		if slices.ContainsFunc(ingRanges, func(r ingRange) bool { return r.containsIngredient(i) }) {
			partA++
		}
	}
	fmt.Printf("Number of fresh ingredients from list: %d\n\n", partA)

	fmt.Println("Part B:")
	partB := 0
	for _, r := range ingRanges {
		partB += (r.max - r.min) + 1
	}
	fmt.Printf("Total number of fresh ingredients: %d\n\n", partB)

}

type ingRange struct {
	min, max int
}

func (r ingRange) intersects(other ingRange) bool {
	return other.max >= r.min && other.min <= r.max
}

func (r ingRange) contains(other ingRange) bool {
	return other.min >= r.min && other.max <= r.max
}

func (r ingRange) containsIngredient(i int) bool {
	return r.min <= i && i <= r.max
}

func newIngRange(s string) ingRange {
	p := strings.Split(s, "-")
	return ingRange{
		min: convert.StrToInt(p[0]),
		max: convert.StrToInt(p[1]),
	}
}

func normaliseIngRanges(ingRanges []ingRange) []ingRange {
	var norm []ingRange
normLoop:
	for i, r := range ingRanges {
		// If intersects with already normalised, skip
		for _, n := range norm {
			if n.contains(r) {
				continue normLoop
			}
		}
		// If intersects with proceeding ranges, merge
		nxtIdx := i + 1
		if nxtIdx < len(ingRanges) {
			for ni := nxtIdx; ni < len(ingRanges); ni++ {
				if r.intersects(ingRanges[ni]) && !r.contains(ingRanges[ni]) {
					r = mergeIngRanges(r, ingRanges[ni])
					ni = nxtIdx - 1 // reset loop, as r might now intersect with ranges we've already checked
					continue
				}
			}

		}
		// Add
		norm = append(norm, r)
	}
	return norm
}

func mergeIngRanges(a, b ingRange) ingRange {
	return ingRange{
		min: min(a.min, b.min),
		max: max(a.max, b.max),
	}
}
