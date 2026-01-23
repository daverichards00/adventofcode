package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/internal/convert"
	"github.com/daverichards00/adventofcode/internal/file"
	"strings"
)

func main() {
	fmt.Println("Day 14")

	input := file.Load("cmd/2019/day14/input.txt")
	for _, line := range input {
		p := strings.Split(line, " => ")
		o := strings.SplitN(p[1], " ", 2)
		r := &reaction{qty: convert.StrToInt(o[0]), req: make(map[string]int)}
		for _, pp := range strings.Split(p[0], ", ") {
			req := strings.Split(pp, " ")
			r.req[req[1]] = convert.StrToInt(req[0])
		}
		chemicals[o[1]] = r
	}

	fmt.Println("Part A:")
	fuelFor1 := requiredOre(1)
	fmt.Printf("Minimum ORE required for 1 FUEL: %d\n\n", fuelFor1)

	fmt.Println("Part B:")
	totalOre := 1000000000000

	// We know the absolute min, binary search for the nearest value to totalOre
	totalFuelMin := totalOre / fuelFor1
	totalFuelMax := totalFuelMin * 2
	for {
		t := totalFuelMin + (totalFuelMax-totalFuelMin)/2
		if t == totalFuelMin {
			break
		}
		if requiredOre(t) > totalOre {
			totalFuelMax = t
		} else {
			totalFuelMin = t
		}
	}

	fmt.Printf("Maximum FUEL with %d ORE: %d\n\n", totalOre, totalFuelMin)
}

type chemical interface {
	needMin(rec int)
	min() int
	reset()
}

var chemicals = map[string]chemical{"ORE": &raw{}}

type raw struct {
	recMin int
}

func (p *raw) needMin(rec int) {
	p.recMin += rec
}

func (p *raw) min() int {
	return p.recMin
}

func (p *raw) reset() {
	p.recMin = 0
}

type reaction struct {
	req           map[string]int
	qty           int
	recMin        int
	recMinSurplus int
}

func (r *reaction) needMin(rec int) {
	need := rec - r.recMinSurplus
	if need <= 0 {
		r.recMinSurplus -= rec
		return
	}
	mul, rem := need/r.qty, need%r.qty
	if rem > 0 {
		mul++
	}
	for k, v := range r.req {
		chemicals[k].needMin(v * mul)
	}
	r.recMin += r.qty * mul
	r.recMinSurplus = (r.qty * mul) - need
}

func (r *reaction) min() int {
	return r.recMin
}

func (r *reaction) reset() {
	r.recMin = 0
	r.recMinSurplus = 0
}

func requiredOre(fuel int) int {
	resetChemicals()
	chemicals["FUEL"].needMin(fuel)
	return chemicals["ORE"].min()
}

func resetChemicals() {
	for i := range chemicals {
		chemicals[i].reset()
	}
}
