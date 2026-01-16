package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/cmd/2019/day05/intcode"
	"github.com/daverichards00/adventofcode/internal/file"
	"github.com/daverichards00/adventofcode/internal/slice"
	"sync"
)

func main() {
	fmt.Println("Day 07")
	input := file.Load("cmd/2019/day07/input.txt")
	program := intcode.NewProgram(input[0])

	fmt.Println("Part A:")
	maxSignal := 0
	for _, combo := range slice.Combinations([]int{0, 1, 2, 3, 4}) {
		p := newPipe(program, 5)
		for i, c := range combo {
			p.inputInstance(i) <- c
		}

		var output int
		var wg sync.WaitGroup
		wg.Go(func() {
			output = <-p.output()
		})

		p.input() <- 0
		wg.Wait()
		maxSignal = max(maxSignal, output)
	}
	fmt.Printf("The maximum signal that can be sent: %d\n\n", maxSignal)

	fmt.Println("Part B:")
	maxFeedbackSignal := 0
	for _, combo := range slice.Combinations([]int{5, 6, 7, 8, 9}) {
		p := newPipe(program, 5)
		for i, c := range combo {
			p.inputInstance(i) <- c
		}

		var output int
		var wg sync.WaitGroup
		wg.Go(func() {
			for {
				output = <-p.output()
				if p.isRunning(0) {
					// feedback into system
					p.input() <- output
					continue
				}
				break
			}
		})

		p.input() <- 0
		wg.Wait()
		maxFeedbackSignal = max(maxFeedbackSignal, output)
	}
	fmt.Printf("The maximum feedback signal that can be sent: %d\n\n", maxFeedbackSignal)
}

type pipe struct {
	ch   []chan int
	exit []*intcode.ExitStatus
	mu   sync.RWMutex
}

func (p *pipe) input() chan<- int {
	return p.ch[0]
}

func (p *pipe) inputInstance(inst int) chan<- int {
	return p.ch[inst]
}

func (p *pipe) output() <-chan int {
	return p.ch[len(p.ch)-1]
}

func (p *pipe) isRunning(inst int) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.exit[inst] == nil
}

func newPipe(program intcode.Program, size int) *pipe {
	p := &pipe{
		ch:   make([]chan int, size+1),
		exit: make([]*intcode.ExitStatus, size),
	}

	for i := 0; i < len(p.ch); i++ {
		p.ch[i] = make(chan int)
	}

	for i := size - 1; i >= 0; i-- {
		go func(i int) {
			exit := intcode.Run(program, p.ch[i], p.ch[i+1])

			e := <-exit

			p.mu.Lock()
			p.exit[i] = &e
			p.mu.Unlock()

			// Close output channel
			close(p.ch[i+1])
			if i == 0 {
				// This is the first instance, close the input channel too
				close(p.ch[i])
			}
		}(i)
	}

	return p
}
