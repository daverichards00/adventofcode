package main

import (
	"errors"
	"fmt"
	"github.com/daverichards00/adventofcode/cmd/2019/day15/intcode"
	"github.com/daverichards00/adventofcode/internal/convert"
	"github.com/daverichards00/adventofcode/internal/file"
	"github.com/daverichards00/adventofcode/internal/spatial"
	"slices"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Day 17")

	input := file.Load("cmd/2019/day17/input.txt")
	program := intcode.NewProgram(input[0])

	output, _ := intcode.RunSync(program, []int{})
	s, r := asciiToSpace(output)

	in := intersections(s)

	fmt.Println("Part A:")
	alignParams := 0
	for _, p := range in {
		alignParams += p.X() * p.Y()
	}
	fmt.Printf("Sum of alignment parameters: %d\n\n", alignParams)

	fmt.Println("Part B:")
	// Find routine
	routines, seq, _ := findPaths(s, r, in)

	// Convert to ascii
	routineInput := routinesToAscii(routines, seq)

	// Video feed
	routineInput = append(routineInput, int('n'), 10)

	// Enable movement and run routines
	program[0] = 2
	output, _ = intcode.RunSync(program, routineInput)

	fmt.Printf("Amount of dust collected: %d\n\n", output[len(output)-1])
}

func asciiToSpace(ascii []int) (*spatial.Space2D[bool], robot) {
	s := spatial.NewSpace2D[bool]()
	var rb robot
	x, y := 0, 0
	for _, a := range ascii {
		if a == 10 {
			x, y = 0, y+1
			continue
		}
		p := spatial.NewPoint2D(x, y)
		r := rune(a)
		switch r {
		case '#':
			s.Set(p, true)
		case '^', 'v', '<', '>':
			// Robot
			rb = robot{p: p}
			switch r {
			case '^':
				rb.v = spatial.North2D
			case 'v':
				rb.v = spatial.South2D
			case '<':
				rb.v = spatial.West2D
			case '>':
				rb.v = spatial.East2D
			}
			// Robot always on scaffold
			s.Set(p, true)
		}
		x++
	}
	return s, rb
}

func intersections(s *spatial.Space2D[bool]) []spatial.Point2D {
	var i []spatial.Point2D
	for _, p := range s.Find(true) {
		if s.Exists(p.Add(spatial.North2D)) && s.Exists(p.Add(spatial.South2D)) && s.Exists(p.Add(spatial.West2D)) && s.Exists(p.Add(spatial.East2D)) {
			i = append(i, p)
		}
	}
	return i
}

type robot struct {
	p spatial.Point2D
	v spatial.Vector2D
}

type path struct {
	instructions []string
	justRotated  bool
	robot        robot
	visited      []spatial.Point2D
}

func findPaths(s *spatial.Space2D[bool], rb robot, intersections []spatial.Point2D) ([][]string, []int, error) {
	paths := []path{
		{robot: rb, visited: []spatial.Point2D{rb.p}},
	}

	for len(paths) > 0 {
		pth := paths[len(paths)-1]
		paths = paths[:len(paths)-1]

		f := pth.robot.p.Add(pth.robot.v)
		if s.Exists(f) && !slices.Contains(pth.visited, f) {
			// Can go forward
			np := f
			for s.Exists(np.Add(pth.robot.v)) && !slices.Contains(intersections, np) {
				np = np.Add(pth.robot.v)
			}

			vd := visited(pth.robot.p, np)
			in := slices.Clone(pth.instructions)
			if len(in) > 0 && !pth.justRotated {
				in[len(in)-1] = strconv.Itoa(convert.StrToInt(in[len(in)-1]) + len(vd))
			} else {
				in = append(in, strconv.Itoa(len(vd)))
			}
			paths = append(paths, path{
				instructions: in,
				justRotated:  false,
				robot:        robot{np, pth.robot.v},
				visited:      append(slices.Clone(pth.visited), vd...),
			})
		}

		if !pth.justRotated {
			// Can we turn left?
			l := pth.robot.v.RotateAntiClockwise90()
			lnp := pth.robot.p.Add(l)
			if s.Exists(lnp) && !slices.Contains(pth.visited, lnp) {
				// Turn left
				paths = append(paths, path{
					instructions: append(slices.Clone(pth.instructions), "L"),
					justRotated:  true,
					robot:        robot{pth.robot.p, l},
					visited:      pth.visited,
				})
			}

			// Can we turn right?
			r := pth.robot.v.RotateClockwise90()
			rnp := pth.robot.p.Add(r)
			if s.Exists(rnp) && !slices.Contains(pth.visited, rnp) {
				// Turn right
				paths = append(paths, path{
					instructions: append(slices.Clone(pth.instructions), "R"),
					justRotated:  true,
					robot:        robot{pth.robot.p, r},
					visited:      pth.visited,
				})
			}

			// Neither? Must be complete
			if !s.Exists(lnp) && !s.Exists(rnp) && !slices.Contains(intersections, pth.robot.p) && visitedAll(s, pth) {
				// Complete
				if r, s, err := findRoutines(pth.instructions, [][]string{}, []int{}); err == nil {
					return r, s, nil
				}
			}
		}
	}

	return nil, nil, errors.New("no paths found")
}

func visited(from, to spatial.Point2D) []spatial.Point2D {
	var v []spatial.Point2D
	u := from.To(to).Unit()
	for from != to {
		from = from.Add(u)
		v = append(v, from)
	}
	return v
}

func visitedAll(s *spatial.Space2D[bool], pth path) bool {
	for _, p := range s.Find(true) {
		if !slices.Contains(pth.visited, p) {
			return false
		}
	}
	return true
}

const (
	minFuncSize  = 4
	routineLimit = 3
)

func findRoutines(remaining []string, routines [][]string, seq []int) ([][]string, []int, error) {
	// Find start
trimRoutineLoop:
	for len(remaining) > 0 {
		for i, r := range routines {
			if len(r) <= len(remaining) && slices.Equal(r, remaining[0:len(r)]) {
				remaining = remaining[len(r):]
				seq = append(seq, i)
				continue trimRoutineLoop
			}
		}
		break
	}
	if len(routines) == routineLimit {
		if len(remaining) > 0 || len(seq) > 10 {
			return routines, seq, errors.New("routine limit reached")
		}
		return routines, seq, nil
	}
	maxRoutineSize := min(len(remaining), 20)
	for !validFunc(remaining[0:maxRoutineSize]) {
		maxRoutineSize--
	}
	for size := maxRoutineSize; size >= minFuncSize; size-- {
		if r, s, err := findRoutines(
			remaining,
			append(slices.Clone(routines), remaining[0:size]),
			slices.Clone(seq),
		); err == nil {
			return r, s, nil
		}
	}
	return routines, seq, errors.New("routines not found")
}

func validFunc(pattern []string) bool {
	return len(funcToString(pattern)) <= 20
}

func funcToString(pattern []string) string {
	return strings.Join(pattern, ",")
}

func routinesToAscii(routines [][]string, seq []int) []int {
	var a []int

	// seq
	var ss []string
	for _, s := range seq {
		switch s {
		case 0:
			ss = append(ss, "A")
		case 1:
			ss = append(ss, "B")
		case 2:
			ss = append(ss, "C")
		}
	}
	a = append(a, strToAscii(strings.Join(ss, ","))...)
	a = append(a, 10)

	for _, r := range routines {
		a = append(a, strToAscii(strings.Join(r, ","))...)
		a = append(a, 10)
	}

	return a
}

func strToAscii(s string) []int {
	var a []int
	for _, r := range s {
		a = append(a, int(r))
	}
	return a
}
