package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/internal/convert"
	"github.com/daverichards00/adventofcode/internal/file"
	"github.com/daverichards00/adventofcode/internal/pathfinder"
	"github.com/daverichards00/adventofcode/internal/slice"
	"strings"
)

func main() {
	fmt.Println("Day 21")
	input := file.Load("cmd/2024/day21/input.txt")

	fmt.Println("Part A:")
	partAKeyPads := []*keypad{&alphaKeypad}
	for i := 0; i < 2; i++ {
		partAKeyPads = append(partAKeyPads, &directionKeypad)
	}
	resetCache(len(partAKeyPads))
	partA := 0
	for _, code := range input {
		partA += shortestInputs(code, partAKeyPads) * convert.StrToInt(strings.Trim(code, "A"))
	}
	fmt.Printf("The sum of complexities: %d\n\n", partA)

	fmt.Println("Part B:")
	partBKeyPads := []*keypad{&alphaKeypad}
	for i := 0; i < 25; i++ {
		partBKeyPads = append(partBKeyPads, &directionKeypad)
	}
	resetCache(len(partBKeyPads))
	partB := 0
	for _, code := range input {
		partB += shortestInputs(code, partBKeyPads) * convert.StrToInt(strings.Trim(code, "A"))
	}
	fmt.Printf("The sum of complexities: %d\n\n", partB)
}

var cache = map[int]map[string]int{}

func resetCache(size int) {
	cache = map[int]map[string]int{}
	for i := 1; i <= size; i++ {
		cache[i] = map[string]int{}
	}
}

func shortestInputs(inputs string, keypads []*keypad) int {
	if cached, ok := cache[len(keypads)][inputs]; ok {
		return cached
	}
	var shortest int
	from := 'A'
	for _, c := range inputs {
		var s int
		for _, p := range keypads[0].paths(from, c) {
			p += "A"
			if len(keypads) > 1 {
				if ss := shortestInputs(p, keypads[1:]); s == 0 || ss < s {
					s = ss
				}
			} else {
				if s == 0 || len(p) < s {
					s = len(p)
				}
			}
		}
		shortest += s
		from = c
	}
	cache[len(keypads)][inputs] = shortest
	return shortest
}

type keypad struct {
	keys      map[rune]map[rune]rune
	pathCache map[string][]string
}

var alphaKeypad = keypad{
	keys: map[rune]map[rune]rune{
		'0': {'^': '2', '>': 'A'},
		'1': {'^': '4', '>': '2'},
		'2': {'^': '5', 'v': '0', '<': '1', '>': '3'},
		'3': {'^': '6', 'v': 'A', '<': '2'},
		'4': {'^': '7', 'v': '1', '>': '5'},
		'5': {'^': '8', 'v': '2', '<': '4', '>': '6'},
		'6': {'^': '9', 'v': '3', '<': '5'},
		'7': {'v': '4', '>': '8'},
		'8': {'v': '5', '<': '7', '>': '9'},
		'9': {'v': '6', '<': '8'},
		'A': {'^': '3', '<': '0'},
	},
	pathCache: map[string][]string{},
}

var directionKeypad = keypad{
	keys: map[rune]map[rune]rune{
		'^': {'v': 'v', '>': 'A'},
		'v': {'^': '^', '<': '<', '>': '>'},
		'<': {'>': 'v'},
		'>': {'^': 'A', '<': 'v'},
		'A': {'v': '>', '<': '^'},
	},
	pathCache: map[string][]string{},
}

type instruction struct {
	keys       string
	directions string
}

func (i instruction) curr() rune {
	return rune(i.keys[len(i.keys)-1])
}

func (i instruction) add(d, k rune) instruction {
	return instruction{fmt.Sprintf("%s%c", i.keys, k), fmt.Sprintf("%s%c", i.directions, d)}
}

func (k *keypad) paths(from, to rune) []string {
	if cached, ok := k.pathCache[fmt.Sprintf("%c%c", from, to)]; ok {
		return cached
	}
	paths, err := pathfinder.Find([]instruction{{keys: string(from)}}, pathfinder.Fn[instruction]{
		Next: func(i instruction) []instruction {
			var next []instruction
			for d, r := range k.keys[i.curr()] {
				if strings.Contains(i.keys, string(r)) {
					continue
				}
				next = append(next, i.add(d, r))
			}
			return next
		},
		Complete: func(i instruction) bool {
			return i.curr() == to
		},
		Less: func(a, b instruction) bool {
			return len(a.keys) < len(b.keys)
		},
	})
	if err != nil {
		panic(err)
	}

	pathInputs := slice.Transform(paths, func(i instruction) string {
		return i.directions
	})
	k.pathCache[fmt.Sprintf("%c%c", from, to)] = pathInputs

	return pathInputs
}
