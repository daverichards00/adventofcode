package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/internal/convert"
	"github.com/daverichards00/adventofcode/internal/file"
)

func main() {
	fmt.Println("Day 01")

	input := file.Load("cmd/2025/day01/input.txt")
	rots := rotations(input)

	fmt.Println("Part A:")
	partA, partB := 0, 0

	t := 50
	for _, r := range rots {
		t = t + r
		if t >= 100 {
			partB += t / 100
		}
		if t < 1 {
			partB += t / -100
			if t-r > 0 {
				partB++
			}
		}
		t = normalise(t, 100)
		if t == 0 {
			partA++
		}
	}

	fmt.Printf("Frequency of 0s: %d\n\n", partA)

	fmt.Println("Part B:")
	fmt.Printf("Frequency of clicks past 0: %d\n\n", partB)
}

func rotations(input []string) []int {
	var rot []int
	for _, line := range input {
		switch line[0] {
		case 'L':
			rot = append(rot, -1*convert.StrToInt(line[1:]))
		case 'R':
			rot = append(rot, convert.StrToInt(line[1:]))
		default:
			panic("invalid rotation")
		}
	}
	return rot
}

func normalise(i, limit int) int {
	i = i % limit
	if i < 0 {
		return i + limit
	}
	return i
}
