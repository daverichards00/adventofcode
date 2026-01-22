package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/cmd/2019/day09/intcode"
	"github.com/daverichards00/adventofcode/internal/file"
	"github.com/daverichards00/adventofcode/internal/spatial"
)

func main() {
	fmt.Println("Day 13")

	input := file.Load("cmd/2019/day13/input.txt")
	program := intcode.NewProgram(input[0])

	fmt.Println("Part A:")
	output, err := intcode.RunSync(program, []int{})
	if err != nil {
		panic(err)
	}
	blockCount := 0
	for i := 2; i < len(output); i += 3 {
		if output[i] == 2 {
			blockCount++
		}
	}
	fmt.Printf("Number of block tiles: %d\n\n", blockCount)

	fmt.Println("Part B:")
	program[0] = 2
	score := play(program)
	fmt.Printf("Final score: %d\n\n", score)
	fmt.Println("To run the game see cmd/2019/day13/game")
}

type game struct {
	tiles    *spatial.Space2D[int]
	score    int
	paddle   spatial.Point2D
	ball     spatial.Point2D
	nextMove int
}

func (g *game) setPaddle(paddle spatial.Point2D) {
	g.paddle = paddle
	g.calcNextMove()
}

func (g *game) setBall(ball spatial.Point2D) {
	g.ball = ball
	g.calcNextMove()
}

func (g *game) calcNextMove() {
	switch {
	case g.ball.X() < g.paddle.X():
		g.nextMove = -1
	case g.ball.X() > g.paddle.X():
		g.nextMove = 1
	default:
		g.nextMove = 0
	}
}

func play(program intcode.Program) int {
	g := &game{tiles: spatial.NewSpace2D[int]()}

	input, output := make(chan int), make(chan int)

	exit := intcode.Run(program, input, output)

	for {
		select {
		case <-exit:
			return g.score
		case x := <-output:
			y := <-output
			t := <-output
			if x == -1 && y == 0 {
				g.score = t
				break
			}
			p := spatial.NewPoint2D(x, y)
			g.tiles.Set(p, t)
			if t == 3 {
				g.setPaddle(p)
			}
			if t == 4 {
				g.setBall(p)
			}
		case input <- g.nextMove:
		}
	}
}
