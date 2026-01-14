package main

import (
	"fmt"
	"slices"
)

var (
	inputMin = 278384
	inputMax = 824795
)

func main() {
	fmt.Println("Day 04")

	fmt.Println("Part A:")
	partA := 0
	for p := inputMin; p <= inputMax; p++ {
		if validPasswordA(p) {
			partA++
		}
	}
	fmt.Printf("Number of valid passwords: %d\n\n", partA)

	fmt.Println("Part B:")
	partB := 0
	for p := inputMin; p <= inputMax; p++ {
		if validPasswordB(p) {
			partB++
		}
	}
	fmt.Printf("Number of valid passwords: %d\n\n", partB)
}

func validPasswordA(p int) bool {
	pd := digits(p)
	return adjacentDigits(pd) && neverDecrease(pd)
}

func validPasswordB(p int) bool {
	pd := digits(p)
	return adjacentDigitsExcl(pd) && neverDecrease(pd)
}

func digits(i int) []int {
	var d []int
	for i > 0 {
		d = append(d, i%10)
		i /= 10
	}
	slices.Reverse(d)
	return d
}

func adjacentDigits(p []int) bool {
	for i := 0; i < len(p)-1; i++ {
		if p[i] == p[i+1] {
			return true
		}
	}
	return false
}

func adjacentDigitsExcl(p []int) bool {
	for i := 0; i < len(p)-1; i++ {
		if p[i] == p[i+1] && (i == len(p)-2 || p[i] != p[i+2]) && (i == 0 || p[i] != p[i-1]) {
			return true
		}
	}
	return false
}

func neverDecrease(p []int) bool {
	for i := 0; i < len(p)-1; i++ {
		if p[i] > p[i+1] {
			return false
		}
	}
	return true
}
