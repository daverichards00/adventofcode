package intcode

import (
	"fmt"
	"github.com/daverichards00/adventofcode/internal/convert"
	"slices"
	"strings"
	"sync"
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

func RunSync(program Program, input []int) ([]int, error) {
	in, out := make(chan int), make(chan int)

	var output []int
	var wg sync.WaitGroup
	wg.Go(func() {
		for o := range out {
			output = append(output, o)
		}
	})

	exit := Run(program, in, out)

	for _, i := range input {
		in <- i
	}

	status := <-exit
	close(in)
	close(out)

	wg.Wait()

	return output, status.Err
}

type ExitStatus struct {
	Err   error
	State []int
	Ptr   int
}

func Run(program Program, input <-chan int, output chan<- int) <-chan ExitStatus {
	exit := make(chan ExitStatus)

	go func() {
		defer close(exit)

		mem := program.Read()
		ptr := 0

		for {
			instr := newInstruction(mem, ptr)

			switch instr.opcode {
			case 1:
				// Add
				mem[instr.params[2]] = mem[instr.params[0]] + mem[instr.params[1]]
				ptr += 4
			case 2:
				// Multiply
				mem[instr.params[2]] = mem[instr.params[0]] * mem[instr.params[1]]
				ptr += 4
			case 3:
				// Read input
				mem[instr.params[0]] = <-input
				ptr += 2
			case 4:
				// Write output
				output <- mem[instr.params[0]]
				ptr += 2
			case 5:
				// Jump if true
				if mem[instr.params[0]] != 0 {
					ptr = mem[instr.params[1]]
				} else {
					ptr += 3
				}
			case 6:
				// Jump if false
				if mem[instr.params[0]] == 0 {
					ptr = mem[instr.params[1]]
				} else {
					ptr += 3
				}
			case 7:
				// Less than
				if mem[instr.params[0]] < mem[instr.params[1]] {
					mem[instr.params[2]] = 1
				} else {
					mem[instr.params[2]] = 0
				}
				ptr += 4
			case 8:
				//Equals
				if mem[instr.params[0]] == mem[instr.params[1]] {
					mem[instr.params[2]] = 1
				} else {
					mem[instr.params[2]] = 0
				}
				ptr += 4
			case 99:
				exit <- ExitStatus{nil, mem, ptr}
				return
			default:
				exit <- ExitStatus{fmt.Errorf("invalid opcode %d", instr.opcode), mem, ptr}
				return
			}
		}
	}()

	return exit
}

type instruction struct {
	opcode int
	params []int
}

func newInstruction(memory []int, ptr int) instruction {
	val := memory[ptr]
	i := instruction{opcode: val % 100}
	if ptr < len(memory) {
		i.params = slices.Clone(memory[ptr+1:])
	}

	var modes []int
	val /= 100
	for val > 0 {
		modes, val = append(modes, val%10), val/10
	}

	for k, v := range modes {
		if v == 1 {
			i.params[k] = ptr + 1 + k
		}
	}

	return i
}
