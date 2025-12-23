package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/internal/convert"
	"github.com/daverichards00/adventofcode/internal/file"
	"github.com/daverichards00/adventofcode/internal/slice"
	"slices"
)

func main() {
	fmt.Println("Day 22")
	input := file.Load("cmd/2024/day22/input.txt")
	inputInts := slice.Transform(input, func(s string) int {
		return convert.StrToInt(s)
	})

	fmt.Println("Part A:")
	partA := 0
	for _, in := range inputInts {
		partA += secretNth(in, 2000)
	}
	fmt.Printf("Sum of each 2000th secret: %d\n\n", partA)

	fmt.Println("Part B:")
	initPrices(inputInts)
	fmt.Printf("Most number of bananas: %d\n\n", most())

}

func secretNth(s, n int) int {
	for i := 0; i < n; i++ {
		s = secret(s)
	}
	return s
}

func secret(s int) int {
	s = (s ^ (s * 64)) % 16777216
	s = (s ^ (s / 32)) % 16777216
	return (s ^ (s * 2048)) % 16777216
}

var prices [][]int
var diffs [][]int

func initPrices(inits []int) {
	prices, diffs = make([][]int, len(inits)), make([][]int, len(inits))
	for i := range inits {
		prices[i], diffs[i] = make([]int, 2001), make([]int, 2000)
		s := inits[i]
		prices[i][0] = s % 10
		for j := 0; j < 2000; j++ {
			s = secret(s)
			prices[i][j+1] = s % 10
			diffs[i][j] = prices[i][j+1] - prices[i][j]
		}
	}
}

var seqLen = 4

func most() int {
	totals := map[int]int{}
	for i := range prices {
		var pastKeys []int
		for j := seqLen; j < len(prices[i]); j++ {
			key := seqKey(diffs[i][j-seqLen : j])
			if slices.Contains(pastKeys, key) {
				continue
			}
			totals[key] = totals[key] + prices[i][j]
			pastKeys = append(pastKeys, key)
		}
	}
	m := 0
	for _, v := range totals {
		m = max(m, v)
	}
	return m
}

func seqKey(seq []int) int {
	k := 0
	for i, mult := len(seq)-1, 1; i >= 0; i, mult = i-1, mult*100 {
		k += mult * (seq[i] + 9)
	}
	return k
}
