package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/internal/file"
	"strings"
)

var devices = map[string][]string{}

func main() {
	fmt.Println("Day 11")

	input := file.Load("cmd/2025/day11/input.txt")
	for _, line := range input {
		d := strings.Split(line, ":")
		outputs := strings.Split(strings.TrimSpace(d[1]), " ")
		devices[d[0]] = outputs
	}

	fmt.Println("Part A:")
	partA := pathCount("you", "out")
	fmt.Printf("Number of possible paths: %d\n\n", partA)

	fmt.Println("Part B:")
	svr2dac := pathCount("svr", "dac")
	dac2fft := pathCount("dac", "fft")
	fft2out := pathCount("fft", "out")
	svr2fft := pathCount("svr", "fft")
	fft2dac := pathCount("fft", "dac")
	dac2out := pathCount("dac", "out")
	partB := (svr2dac * dac2fft * fft2out) + (svr2fft * fft2dac * dac2out)
	fmt.Printf("Number of possible paths: %d\n\n", partB)
}

var pathCache = make(cache)

func pathCount(from, to string) int {
	if total, ok := pathCache.get(from, to); ok {
		return total
	}
	total := 0
	for _, out := range devices[from] {
		if out == to {
			total++
			continue
		}
		total += pathCount(out, to)
	}
	pathCache.set(from, to, total)
	return total
}

type cache map[string]map[string]int

func (c cache) get(from, to string) (int, bool) {
	f, ok := c[from]
	if !ok {
		return 0, false
	}
	t, ok := f[to]
	return t, ok
}

func (c cache) set(from, to string, value int) {
	if _, ok := c[from]; !ok {
		c[from] = make(map[string]int)
	}
	c[from][to] = value
}
