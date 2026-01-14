package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/cmd/2019/day02/intcode"
	"github.com/daverichards00/adventofcode/internal/file"
)

func main() {
	fmt.Println("Day 02")
	input := file.Load("cmd/2019/day02/input.txt")

	program := intcode.NewProgram(input[0])

	fmt.Println("Part A:")
	program[1] = 12
	program[2] = 2
	stateA, _ := intcode.Run(program)
	fmt.Printf("Value at position 0: %d\n\n", stateA[0])

	fmt.Println("Part B:")
	for noun := 0; noun < 100; noun++ {
		for verb := 0; verb < 100; verb++ {
			program[1] = noun
			program[2] = verb
			state, err := intcode.Run(program)
			if err != nil {
				panic(err)
			}
			if state[0] == 19690720 {
				fmt.Printf("Noun/verb that produces 19690720: %d\n\n", noun*100+verb)
				return
			}
		}
	}
}
