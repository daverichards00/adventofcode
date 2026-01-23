package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/cmd/2019/day15/intcode"
	"github.com/daverichards00/adventofcode/internal/file"
	"github.com/daverichards00/adventofcode/internal/pathfinder"
	"github.com/daverichards00/adventofcode/internal/spatial"
	"slices"
)

const (
	wall = iota
	space
	complete
)

var (
	program intcode.Program

	inputDirections = map[spatial.Vector2D]int{
		spatial.North2D: 1,
		spatial.South2D: 2,
		spatial.West2D:  3,
		spatial.East2D:  4,
	}
)

func main() {
	fmt.Println("Day 15")

	input := file.Load("cmd/2019/day15/input.txt")
	program = intcode.NewProgram(input[0])

	fmt.Println("Part A:")
	s, c := discover()
	fmt.Printf("Fewest movements to oxygen system: %d\n\n", len(c)-1)

	fmt.Println("Part B:")
	o := oxygenate(s, s.Find(complete)[0])
	last := 0
	for _, t := range o {
		last = max(last, t)
	}
	fmt.Printf("Time to fill area: %d\n\n", last)
}

func valid(p []spatial.Point2D) (int, error) {
	input, output := make(chan int), make(chan int)

	exit, sigterm := intcode.Run(program, input, output)
	defer func() {
		sigterm <- struct{}{}
		close(input)
		close(output)
	}()

	robot := p[0]
	for i := 1; i < len(p); i++ {
		move := robot.To(p[i])
		input <- inputDirections[move]
		select {
		case e := <-exit:
			return 0, e.Err
		case out := <-output:
			switch out {
			case wall:
				return wall, nil
			case space:
				robot = p[i]
			case complete:
				return complete, nil
			default:
				return 0, fmt.Errorf("unknown output: %d", out)
			}
		}
	}
	return space, nil
}

var dirs = []spatial.Vector2D{spatial.North2D, spatial.South2D, spatial.West2D, spatial.East2D}

func discover() (*spatial.Space2D[int], []spatial.Point2D) {
	s := spatial.NewSpace2D[int]()
	var c []spatial.Point2D
	paths := [][]spatial.Point2D{{spatial.NewPoint2D(0, 0)}}

	_, _ = pathfinder.Find(paths, pathfinder.Fn[[]spatial.Point2D]{
		Next: func(path []spatial.Point2D) [][]spatial.Point2D {
			var next [][]spatial.Point2D
			last := path[len(path)-1]
			for _, dir := range dirs {
				np := last.Add(dir)
				if s.Exists(np) {
					continue
				}
				n := append(slices.Clone(path), np)
				v, _ := valid(n)
				s.Set(np, v)

				if v == wall {
					continue
				}
				if v == complete && len(c) == 0 {
					c = n
				}
				next = append(next, n)
			}
			return next
		},
		Complete: func(path []spatial.Point2D) bool {
			return false
		},
		Less: func(a, b []spatial.Point2D) bool {
			return len(a) < len(b)
		},
	})

	return s, c
}

func oxygenate(s *spatial.Space2D[int], start spatial.Point2D) map[spatial.Point2D]int {
	o := map[spatial.Point2D]int{start: 0}
	paths := [][]spatial.Point2D{{start}}

	_, _ = pathfinder.Find(paths, pathfinder.Fn[[]spatial.Point2D]{
		Next: func(path []spatial.Point2D) [][]spatial.Point2D {
			var next [][]spatial.Point2D
			last := path[len(path)-1]
			for _, dir := range dirs {
				np := last.Add(dir)
				if _, ok := o[np]; ok {
					continue
				}
				if s.GetOrDefault(np, wall) == wall {
					continue
				}
				n := append(slices.Clone(path), np)

				o[np] = len(n) - 1

				next = append(next, n)
			}
			return next
		},
		Complete: func(path []spatial.Point2D) bool {
			return false
		},
		Less: func(a, b []spatial.Point2D) bool {
			return len(a) < len(b)
		},
	})
	return o
}
