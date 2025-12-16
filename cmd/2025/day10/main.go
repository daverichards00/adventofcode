package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/internal/convert"
	"github.com/daverichards00/adventofcode/internal/file"
	"github.com/daverichards00/adventofcode/internal/numrange"
	"github.com/daverichards00/adventofcode/internal/pathfinder"
	"github.com/daverichards00/adventofcode/internal/slice"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"
)

var (
	inputRegExp = regexp.MustCompile(`^\[([.#]+)]((?: \([0-9]+(?:,[0-9]+)*\))+) \{([0-9]+(?:,[0-9]+)*)}$`)
)

func main() {
	fmt.Println("Day 10")

	var machines []machine

	input := file.Load("cmd/2025/day10/input.txt")
	for _, line := range input {
		m := inputRegExp.FindStringSubmatch(line)

		machines = append(machines, machine{
			targetState:   newTargetState(m[1]),
			buttons:       newButtons(m[2]),
			targetJoltage: newTargetJoltage(m[3]),
		})
	}

	start := time.Now()

	fmt.Println("Part A:")
	partA := 0
	for _, m := range machines {
		s := shortestStateSequence(m)
		partA += s.len()
	}
	fmt.Printf("The fewest number of button presses: %d\n\n", partA)

	// Part B takes about 5m 30s to complete
	fmt.Println("Part B:")
	partB := 0
	for _, m := range machines {
		partB += shortestJoltageSequence(m)
	}
	fmt.Printf("The fewest number of button presses: %d\n\n", partB)

	fmt.Printf("Duration: %s\n\n", time.Since(start).String())
}

type machine struct {
	targetState   state
	targetJoltage joltage
	buttons       []button
}

type state []bool

func (s state) equal(other state) bool {
	if len(s) != len(other) {
		return false
	}
	for i := range s {
		if s[i] != other[i] {
			return false
		}
	}
	return true
}

func (s state) string() string {
	r := ""
	for i := range s {
		if s[i] {
			r += "#"
		} else {
			r += "."
		}
	}
	return r
}

func (s state) pressButton(b button) state {
	n := slices.Clone(s)
	for _, idx := range b {
		n[idx] = !n[idx]
	}
	return n
}

func newState(size int) state {
	return make(state, size)
}

func newTargetState(s string) state {
	n := make(state, len(s))
	for i, c := range s {
		n[i] = c == '#'
	}
	return n
}

type button []int

func newButtons(s string) []button {
	var buttons []button
	for _, p := range strings.Split(strings.TrimSpace(s), " ") {
		var b button
		for _, pp := range strings.Split(strings.Trim(p, "()"), ",") {
			b = append(b, convert.StrToInt(pp))
		}
		buttons = append(buttons, b)
	}
	return buttons
}

type joltage []int

func newJoltage(size int) joltage {
	return make(joltage, size)
}

func newTargetJoltage(s string) joltage {
	var n joltage
	for _, p := range strings.Split(s, ",") {
		n = append(n, convert.StrToInt(p))
	}
	return n
}

type sequence struct {
	state   state
	joltage joltage
	buttons []button
}

func (s sequence) next(b button) sequence {
	return sequence{
		state:   s.state.pressButton(b),
		buttons: append(slices.Clone(s.buttons), b),
	}
}

func (s sequence) len() int {
	return len(s.buttons)
}

func newSequence(s, j int) sequence {
	return sequence{
		state:   newState(s),
		joltage: newJoltage(j),
	}
}

func shortestStateSequence(m machine) sequence {
	init := []sequence{newSequence(len(m.targetState), len(m.targetJoltage))}
	cache := map[string]bool{}

	shortest, err := pathfinder.Find(init, pathfinder.Fn[sequence]{
		Next: func(s sequence) []sequence {
			var n []sequence
			for _, b := range m.buttons {
				nn := s.next(b)
				if _, ok := cache[nn.state.string()]; ok {
					continue
				}
				cache[nn.state.string()] = true
				n = append(n, nn)
			}
			return n
		},
		Complete: func(s sequence) bool {
			return s.state.equal(m.targetState)
		},
		Less: func(a, b sequence) bool {
			return a.len() < b.len()
		},
	})
	if err != nil {
		panic(err)
	}
	return shortest[0]
}

type joltageCounter struct {
	target     int
	buttonIdxs []int
}

func shortestJoltageSequence(m machine) int {
	// Sort buttons in order of how many counters they affect
	buttons := slices.Clone(m.buttons)
	sort.Slice(buttons, func(i, j int) bool { return len(buttons[i]) > len(buttons[j]) })

	// Init ranges of possible button presses
	pressRanges := make([]numrange.NumRange[int], len(buttons))
	initMaxPress := slices.Max(m.targetJoltage)
	for i := range pressRanges {
		pressRanges[i] = numrange.New(0, initMaxPress)
	}

	counters := make([]joltageCounter, len(m.targetJoltage))
	for i := range counters {
		counters[i] = joltageCounter{
			target: m.targetJoltage[i],
			buttonIdxs: slice.IndexAllFunc(buttons, func(b button) bool {
				return slices.Contains(b, i)
			}),
		}
	}

	pressed, ok := balanceRanges(pressRanges, counters, 0)
	if !ok {
		panic("no possible solutions found")
	}

	total := 0
	for _, r := range pressed {
		total += r.Min()
	}
	return total
}

func balanceRanges(pressRanges []numrange.NumRange[int], counters []joltageCounter, beat int) ([]numrange.NumRange[int], bool) {
	for {
		before := slices.Clone(pressRanges)

		for _, c := range counters {
			for _, bIndx := range c.buttonIdxs {
				othersTotalMin, othersTotalMax := 0, 0
				for _, otherIdx := range slice.Filter(c.buttonIdxs, func(i int) bool { return i != bIndx }) {
					othersTotalMin += pressRanges[otherIdx].Min()
					othersTotalMax += min(pressRanges[otherIdx].Max(), c.target)
				}
				intersect, ok := pressRanges[bIndx].Intersect(numrange.New(c.target-othersTotalMax, c.target-othersTotalMin))
				if !ok {
					return pressRanges, false
				}
				pressRanges[bIndx] = intersect
			}
		}
		if slices.Equal(pressRanges, before) {
			// Nothing changed this loop, no more can be done
			break
		}
	}
	if !slices.ContainsFunc(pressRanges, func(r numrange.NumRange[int]) bool {
		return r.Min() != r.Max()
	}) {
		// All balanced
		return pressRanges, true
	}
	var best []numrange.NumRange[int]
	for i := 0; i < len(pressRanges); i++ {
		if pressRanges[i].Min() == pressRanges[i].Max() {
			continue
		}
		for j := pressRanges[i].Min(); j <= pressRanges[i].Max(); j++ {
			nextRanges := slices.Clone(pressRanges)
			nextRanges[i] = numrange.New(j, j)
			if beat > 0 && minTotal(nextRanges) >= beat {
				break
			}
			nextRanges, ok := balanceRanges(slices.Clone(nextRanges), counters, beat)
			if ok {
				m := minTotal(nextRanges)
				if beat > 0 && m >= beat {
					continue
				}
				beat = m
				best = slices.Clone(nextRanges)
			}
		}
		break
	}
	if len(best) == 0 {
		// Couldn't find a valid value for i
		return pressRanges, false
	}
	return best, true
}

func minTotal(ranges []numrange.NumRange[int]) int {
	m := 0
	for _, r := range ranges {
		m += r.Min()
	}
	return m
}
