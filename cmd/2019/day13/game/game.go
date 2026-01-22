package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/cmd/2019/day09/intcode"
	"github.com/daverichards00/adventofcode/internal/file"
	"github.com/daverichards00/adventofcode/internal/spatial"
	"github.com/nsf/termbox-go"
	"time"
)

const (
	gameTick        = 50 * time.Millisecond
	backgroundColor = termbox.ColorBlack
	textColor       = termbox.ColorWhite
)

const (
	stateInit int = iota
	stateActive
	stateOver
)

func main() {
	input := file.Load("cmd/2019/day13/input.txt")
	program := intcode.NewProgram(input[0])

	program[0] = 2

	score := play(program)
	fmt.Printf("Final score: %d\n\n", score)
}

type tile struct {
	rune rune
	fg   termbox.Attribute
	bg   termbox.Attribute
}

var tiles = []tile{
	{' ', termbox.ColorBlack, termbox.ColorBlack},
	{' ', termbox.ColorBlack, termbox.ColorWhite},
	{'#', termbox.ColorBlue, termbox.ColorBlack},
	{'=', termbox.ColorGreen, termbox.ColorBlack},
	{'0', termbox.ColorWhite, termbox.ColorBlack},
}

type game struct {
	tiles  *spatial.Space2D[tile]
	score  int
	paddle spatial.Point2D
	ball   spatial.Point2D
	state  int
}

func (g *game) update(x, y, t int) {
	if x == -1 && y == 0 {
		g.score = t
		return
	}

	p := spatial.NewPoint2D(x, y)
	g.tiles.Set(p, tiles[t])
	if t == 3 {
		g.paddle = p
	}
	if t == 4 {
		g.ball = p
	}
}

func (g *game) render() {
	termbox.Clear(backgroundColor, backgroundColor)

	tbPrintString(1, 1, textColor, backgroundColor, fmt.Sprintf("Score: %d", g.score))

	for p, t := range g.tiles.GetAll() {
		tbPrintString(p.X()+1, p.Y()+3, t.fg, t.bg, string(t.rune))
	}

	if g.state == stateInit {
		tbPrintBox(8, 11, 33, 13, backgroundColor)
		tbPrintString(10, 12, textColor, backgroundColor, "Press [enter] to start")
	}
	if g.state == stateOver {
		tbPrintBox(8, 11, 34, 13, backgroundColor)
		tbPrintString(10, 12, textColor, backgroundColor, "Complete! [esc] to exit")
	}

	termbox.Flush()
}

func play(program intcode.Program) int {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	g := &game{
		tiles: spatial.NewSpace2D[tile](),
		state: stateInit,
	}

	input, output := make(chan int), make(chan int)
	exit := intcode.Run(program, input, output)

	ticker := time.NewTicker(gameTick)
	defer ticker.Stop()

	for {
		select {
		case <-exit:
			g.state = stateOver
		case event := <-eventQueue:
			if event.Type == termbox.EventKey {
				// Can handle keyboard input here
				switch {
				case event.Key == termbox.KeyEnter:
					if g.state == stateInit {
						g.state = stateActive
					}
				case event.Key == termbox.KeyEsc:
					if g.state == stateOver {
						return g.score
					}
				default:
				}
			}
		case x := <-output:
			// Update State from program
			y := <-output
			t := <-output
			g.update(x, y, t)
		case <-ticker.C:
			// Render
			g.render()
			if g.state == stateActive {
				// AI
				switch {
				case g.ball.X() < g.paddle.X():
					input <- -1
				case g.ball.X() > g.paddle.X():
					input <- 1
				default:
					input <- 0
				}
			}
		}
	}
}

func tbPrintString(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func tbPrintBox(x1, y1, x2, y2 int, col termbox.Attribute) {
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			termbox.SetCell(x, y, ' ', col, col)
		}
	}
}
