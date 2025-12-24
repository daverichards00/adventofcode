package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/internal/file"
	"github.com/daverichards00/adventofcode/internal/slice"
	"maps"
	"slices"
	"sort"
	"strings"
)

var network = map[string][]string{}

func main() {
	fmt.Println("Day 23")

	input := file.Load("cmd/2024/day23/input.txt")
	for _, line := range input {
		p := strings.Split(line, "-")
		network[p[0]] = append(network[p[0]], p[1])
		network[p[1]] = append(network[p[1]], p[0])
	}
	for _, n := range network {
		sort.Strings(n)
	}

	computers := slices.Collect(maps.Keys(network))
	sort.Strings(computers)

	fmt.Println("Part A:")
	var setsOf3 [][]string
	for _, c := range computers {
		n := network[c]
		for i := 0; i < len(n)-1; i++ {
			if n[i] < c {
				continue
			}
			for j := i + 1; j < len(n); j++ {
				if slices.Contains(network[n[i]], n[j]) && (c[0] == 't' || n[i][0] == 't' || n[j][0] == 't') {
					setsOf3 = append(setsOf3, []string{c, n[i], n[j]})
				}
			}
		}
	}
	fmt.Printf("Number of networks of 3 that contain a computer beginning with 't': %d\n\n", len(setsOf3))

	fmt.Println("Part B:")
	var sets [][]string
	for _, c := range computers {
		sets = append(sets, expand([]string{c})...)
	}
	sort.Slice(sets, func(i, j int) bool { return len(sets[i]) > len(sets[j]) })

	fmt.Printf("Password to join largest LAN: %s\n\n", strings.Join(sets[0], ","))
}

func expand(computers []string) [][]string {
	intersect := network[computers[0]]
	for i := 1; i < len(computers); i++ {
		intersect = slice.Intersect(intersect, network[computers[i]])
	}
	if len(intersect) == 0 {
		return [][]string{computers}
	}
	var expanded [][]string
	for i := range intersect {
		if slices.ContainsFunc(expanded, func(e []string) bool { return slices.Contains(e, intersect[i]) }) {
			continue
		}
		expanded = append(expanded, expand(append(slices.Clone(computers), intersect[i]))...)
	}
	return expanded
}
