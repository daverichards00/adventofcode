package intcode

import (
	"errors"
	"github.com/daverichards00/adventofcode/internal/convert"
	"slices"
	"strings"
)

type Program []int

func (p Program) Read() []int {
	return slices.Clone(p)
}

func NewProgram(raw string) Program {
	var n Program
	for _, s := range strings.Split(raw, ",") {
		n = append(n, convert.StrToInt(s))
	}
	return n
}

func Run(program Program) ([]int, error) {
	memory := program.Read()
	ptr := 0

	for {
		switch memory[ptr] {
		case 1:
			memory[memory[ptr+3]] = memory[memory[ptr+1]] + memory[memory[ptr+2]]
			ptr += 4
		case 2:
			memory[memory[ptr+3]] = memory[memory[ptr+1]] * memory[memory[ptr+2]]
			ptr += 4
		case 99:
			return memory, nil
		default:
			return memory, errors.New("invalid opcode")
		}
	}
}
