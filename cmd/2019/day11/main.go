package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/cmd/2019/day09/intcode"
	"github.com/daverichards00/adventofcode/internal/file"
	"github.com/daverichards00/adventofcode/internal/spatial"
	"sync"
)

func main() {
	fmt.Println("Day 11")

	input := file.Load("cmd/2019/day11/input.txt")
	program := intcode.NewProgram(input[0])

	fmt.Println("Part A:")
	hull := paint(program, 0)
	fmt.Printf("Number of panels painted: %d\n\n", len(hull.GetAll()))

	fmt.Println("Part B:")
	hull = paint(program, 1)
	fmt.Println(hull.String(func(p spatial.Point2D, col int) rune {
		if col == 1 {
			return '#'
		}
		return ' '
	}, ' '))
}

func paint(program intcode.Program, start int) *spatial.Space2D[int] {
	hull := spatial.NewSpace2D[int]()

	robot := struct {
		p spatial.Point2D
		v spatial.Vector2D
	}{
		p: spatial.NewPoint2D(0, 0),
		v: spatial.North2D,
	}

	input, output := make(chan int), make(chan int)
	defer close(input)
	defer close(output)

	exit := intcode.Run(program, input, output)

	var wg sync.WaitGroup
	wg.Go(func() {
		for {
			select {
			case <-exit:
				return
			case col := <-output:
				hull.Set(robot.p, col)
				switch <-output {
				case 0:
					// Anti-Clockwise
					robot.v = robot.v.RotateAntiClockwise90()
				case 1:
					// Clockwise
					robot.v = robot.v.RotateClockwise90()
				default:
					panic("unexpected output")
				}
				robot.p = robot.p.Add(robot.v)

				select {
				case <-exit:
					return
				case input <- hull.GetOrDefault(robot.p, 0):
				}
			}
		}
	})

	input <- start

	wg.Wait()
	return hull
}
