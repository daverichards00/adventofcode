package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/internal/file"
	"slices"
	"strings"
)

func main() {
	fmt.Println("Day 06")
	input := file.Load("cmd/2019/day06/input.txt")
	for _, line := range input {
		p := strings.Split(line, ")")
		orbitMap[p[1]] = object{direct: p[0]}
	}

	fmt.Println("Part A:")
	partA := 0
	for _, o := range orbitMap {
		partA += len(o.orbits())
	}
	fmt.Printf("Total number of orbits: %d\n\n", partA)

	fmt.Println("Part B:")
	yo := orbitMap["YOU"].orbits()
	for k, v := range orbitMap["SAN"].orbits() {
		if i := slices.Index(yo, v); i != -1 {
			fmt.Printf("Min transfers required: %d\n\n", k+i)
			break
		}
	}
}

var orbitMap = map[string]object{"COM": {}}

type object struct {
	direct   string
	indirect []string
}

func (obj object) orbits() []string {
	if obj.direct == "" {
		return []string{}
	}
	if obj.indirect == nil {
		obj.indirect = orbitMap[obj.direct].orbits()
	}
	return append([]string{obj.direct}, obj.indirect...)
}
