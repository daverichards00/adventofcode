package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/cmd/2019/day05/intcode"
	"github.com/daverichards00/adventofcode/internal/file"
)

func main() {
	fmt.Println("Day 05")
	input := file.Load("cmd/2019/day05/input.txt")

	program := intcode.NewProgram(input[0])

	fmt.Println("Part A:")
	output, _ := intcode.RunSync(program, []int{1})
	fmt.Printf("Diagnostic code for input 1: %d\n\n", output[len(output)-1])

	fmt.Println("Part B:")
	output, _ = intcode.RunSync(program, []int{5})
	fmt.Printf("Diagnostic code for input 5: %d\n\n", output[0])
}
