package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/internal/convert"
	"github.com/daverichards00/adventofcode/internal/file"
	"github.com/daverichards00/adventofcode/internal/maths"
	"github.com/daverichards00/adventofcode/internal/slice"
	"slices"
	"sort"
	"strings"
)

var wires = map[string]outputter{}

func main() {
	fmt.Println("Day 24")

	input := file.Load("cmd/2024/day24/input.txt")
	idx := slices.Index(input, "")
	for _, line := range input[0:idx] {
		p := strings.Split(line, ": ")
		wires[p[0]] = &constant{value: p[1] == "1"}
	}
	for _, line := range input[idx+1:] {
		p := strings.Split(line, " ")
		switch p[1] {
		case "AND":
			wires[p[4]] = &andGate{a: p[0], b: p[2]}
		case "OR":
			wires[p[4]] = &orGate{a: p[0], b: p[2]}
		case "XOR":
			wires[p[4]] = &xorGate{a: p[0], b: p[2]}
		}
	}

	fmt.Println("Part A:")
	partA := 0
	for k, w := range wires {
		if k[0] == 'z' && w.output() {
			partA += maths.PowInt(2, convert.StrToInt(k[1:]))
		}
	}
	fmt.Printf("Decimal output of z wires: %d\n\n", partA)

	fmt.Println("Part B:")
	errOut := findGateOutputErrors(input[idx+1:])

	fmt.Printf("Wires that need to be swapped: %s\n\n", strings.Join(errOut, ","))
}

type outputter interface {
	output() bool
}

type constant struct {
	value bool
}

func (s constant) output() bool {
	return s.value
}

type andGate struct {
	a, b string
}

func (g *andGate) output() bool {
	return wires[g.a].output() && wires[g.b].output()
}

type orGate struct {
	a, b string
}

func (g *orGate) output() bool {
	return wires[g.a].output() || wires[g.b].output()
}

type xorGate struct {
	a, b string
}

func (g *xorGate) output() bool {
	return wires[g.a].output() != wires[g.b].output()
}

func findGateOutputErrors(gatesRaw []string) []string {
	var errOut []string

	var gates [][]string
	for _, line := range gatesRaw {
		gates = append(gates, strings.Split(line, " "))
	}

	carry, _ := findOutput(gates, "x00", "AND", "y00")
	for i := 1; i <= 44; i++ {
		x, y, z := fmt.Sprintf("x%02d", i), fmt.Sprintf("y%02d", i), fmt.Sprintf("z%02d", i)
		xyXor, _ := findOutput(gates, x, "XOR", y)

		zOut, ok := findOutput(gates, xyXor, "XOR", carry)
		if !ok {
			inA, inB, _ := findInputs(gates, "XOR", z)
			if xyXor == inA || xyXor == inB {
				if xyXor == inA {
					errOut = append(errOut, carry, inB)
					gates, carry, _ = swapOutputs(gates, carry, inB)
				} else {
					errOut = append(errOut, carry, inA)
					gates, carry, _ = swapOutputs(gates, carry, inA)
				}
			} else if carry == inA || carry == inB {
				if carry == inA {
					errOut = append(errOut, xyXor, inB)
					gates, xyXor, _ = swapOutputs(gates, xyXor, inB)
				} else {
					errOut = append(errOut, xyXor, inA)
					gates, xyXor, _ = swapOutputs(gates, xyXor, inA)
				}
			}
		} else if zOut != z {
			errOut = append(errOut, zOut, z)
			gates, _, zOut = swapOutputs(gates, z, zOut)
		}

		carryAnd, _ := findOutput(gates, xyXor, "AND", carry)

		xyAnd, _ := findOutput(gates, x, "AND", y)
		carry, _ = findOutput(gates, xyAnd, "OR", carryAnd)
	}

	errOut = slice.Unique(errOut)
	sort.Strings(errOut)
	return errOut
}

func findOutput(gates [][]string, a, op, b string) (string, bool) {
	for _, gate := range gates {
		if gate[1] == op && ((gate[0] == a && gate[2] == b) || (gate[0] == b && gate[2] == a)) {
			return gate[4], true
		}
	}
	return "", false
}

func findInputs(gates [][]string, op, out string) (string, string, bool) {
	for _, gate := range gates {
		if gate[1] == op && gate[4] == out {
			return gate[0], gate[2], true
		}
	}
	return "", "", false
}

func swapOutputs(gates [][]string, a, b string) ([][]string, string, string) {
	var n [][]string
	for _, gate := range gates {
		if gate[4] == a {
			gate[4] = b
		} else if gate[4] == b {
			gate[4] = a
		}
		n = append(n, gate)
	}
	return n, b, a
}
