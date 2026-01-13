package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/internal/convert"
	"github.com/daverichards00/adventofcode/internal/file"
	"github.com/daverichards00/adventofcode/internal/slice"
)

func main() {
	fmt.Println("Day 01")
	input := file.Load("cmd/2019/day01/input.txt")

	fmt.Println("Part A:")
	fuel := slice.Transform(input, func(s string) int {
		return calcFuel(convert.StrToInt(s))
	})
	fmt.Printf("Sum of fuel requirements: %d\n\n", slice.Sum(fuel))

	fmt.Println("Part B:")
	fuelRecurrsive := slice.Transform(input, func(s string) int {
		return calcFuelRecursive(convert.StrToInt(s))
	})
	fmt.Printf("Sum of fuel requirements, allowing for the mass of the fuel: %d\n\n", slice.Sum(fuelRecurrsive))
}

func calcFuel(mass int) int {
	return mass/3 - 2
}

func calcFuelRecursive(mass int) int {
	fuel := calcFuel(mass)
	total := fuel
	for {
		fuel = calcFuel(fuel)
		if fuel < 1 {
			return total
		}
		total += fuel
	}
}
