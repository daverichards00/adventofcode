package intcode

import (
	"fmt"
	"github.com/daverichards00/adventofcode/internal/convert"
	"slices"
	"strings"
	"sync"
)

const maxInstructionParams = 3

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

		mem := newMemory(program.Read())
		ptr, relBase := 0, 0

		for {
			instr := newInstruction(mem, ptr, relBase)

			switch instr.opcode {
			case 1:
				// Add
				mem.write(instr.params[2], mem.read(instr.params[0])+mem.read(instr.params[1]))
				ptr += 4
			case 2:
				// Multiply
				mem.write(instr.params[2], mem.read(instr.params[0])*mem.read(instr.params[1]))
				ptr += 4
			case 3:
				// Read input
				mem.write(instr.params[0], <-input)
				ptr += 2
			case 4:
				// Write output
				output <- mem.read(instr.params[0])
				ptr += 2
			case 5:
				// Jump if true
				if mem.read(instr.params[0]) != 0 {
					ptr = mem.read(instr.params[1])
				} else {
					ptr += 3
				}
			case 6:
				// Jump if false
				if mem.read(instr.params[0]) == 0 {
					ptr = mem.read(instr.params[1])
				} else {
					ptr += 3
				}
			case 7:
				// Less than
				if mem.read(instr.params[0]) < mem.read(instr.params[1]) {
					mem.write(instr.params[2], 1)
				} else {
					mem.write(instr.params[2], 0)
				}
				ptr += 4
			case 8:
				// Equals
				if mem.read(instr.params[0]) == mem.read(instr.params[1]) {
					mem.write(instr.params[2], 1)
				} else {
					mem.write(instr.params[2], 0)
				}
				ptr += 4
			case 9:
				// Adjust relative base
				relBase += mem.read(instr.params[0])
				ptr += 2
			case 99:
				// Exit
				exit <- ExitStatus{nil, mem.mem, ptr}
				return
			default:
				// Unknown opcode
				exit <- ExitStatus{fmt.Errorf("invalid opcode %d", instr.opcode), mem.mem, ptr}
				return
			}
		}
	}()

	return exit
}

type memory struct {
	mem []int
}

func (m *memory) read(addr int) int {
	if addr >= len(m.mem) {
		return 0
	}
	return m.mem[addr]
}

func (m *memory) write(addr, val int) {
	if addr >= len(m.mem) {
		// Expand available memory
		m.mem = append(m.mem, make([]int, addr-len(m.mem)+1)...)
	}
	m.mem[addr] = val
}

func newMemory(init []int) *memory {
	return &memory{init}
}

type instruction struct {
	opcode int
	params []int
}

func newInstruction(mem *memory, ptr, relBase int) instruction {
	val := mem.read(ptr)
	i := instruction{opcode: val % 100, params: make([]int, maxInstructionParams)}
	for p := range i.params {
		i.params[p] = mem.read(ptr + 1 + p)
	}

	var modes []int
	val /= 100
	for val > 0 {
		modes, val = append(modes, val%10), val/10
	}

	// i.params should hold the addr of each param
	for k, v := range modes {
		// 0: Position mode: Value already contains the param addr
		switch v {
		case 1:
			// Immediate mode: This index is the addr itself
			i.params[k] = ptr + 1 + k
		case 2:
			// Relative mode: Similar to position mode, but offset by the relative base
			i.params[k] = relBase + i.params[k]
		}
	}

	return i
}
