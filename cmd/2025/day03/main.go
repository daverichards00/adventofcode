package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/internal/convert"
	"github.com/daverichards00/adventofcode/internal/file"
	"github.com/daverichards00/adventofcode/internal/maths"
)

func main() {
	fmt.Println("Day 03")

	input := file.Load("cmd/2025/day03/input.txt")
	banks := make([][]int, len(input))
	for i := range input {
		banks[i] = convert.StrToDigits(input[i])
	}

	partA, partB := 0, 0
	for _, bank := range banks {
		partA += largestJoltage(bank, 2)
		partB += largestJoltage(bank, 12)
	}

	fmt.Println("Part A:")
	fmt.Printf("Sum of maximum joltages: %d\n\n", partA)

	fmt.Println("Part B:")
	fmt.Printf("Sum of maximum joltages: %d\n\n", partB)
}

func largestJoltage(bank []int, size int) int {
	idx := largestIntIndex(bank[0 : len(bank)-(size-1)])
	if size > 1 {
		return bank[idx]*maths.PowInt(10, size-1) + largestJoltage(bank[idx+1:], size-1)
	}
	return bank[idx]
}

func largestIntIndex(i []int) int {
	idx, val := 0, i[0]
	for j := 1; j < len(i); j++ {
		if i[j] > val {
			idx, val = j, i[j]
		}
	}
	return idx
}
