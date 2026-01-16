package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/cmd/2019/day09/intcode"
	"github.com/daverichards00/adventofcode/internal/file"
)

func main() {
	fmt.Println("Day 09")
	input := file.Load("cmd/2019/day09/input.txt")
	program := intcode.NewProgram(input[0])

	fmt.Println("Part A:")
	output, _ := intcode.RunSync(program, []int{1})
	fmt.Printf("BOOST keycode: %d\n\n", output[0])

	fmt.Println("Part B:")
	output, _ = intcode.RunSync(program, []int{2})
	fmt.Printf("Coordinates of the distress signal: %d\n\n", output[0])
}
