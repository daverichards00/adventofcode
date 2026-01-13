package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/internal/file"
	"github.com/daverichards00/adventofcode/internal/slice"
)

func main() {
	fmt.Println("Day 25")

	var locks [][]int
	var keys [][]int

	input := file.Load("cmd/2024/day25/input.txt")
	inputSplit := slice.Split(input, "")
	for _, is := range inputSplit {
		keyLock, isKey := parseKeyLock(is)
		if isKey {
			keys = append(keys, keyLock)
		} else {
			locks = append(locks, keyLock)
		}
	}

	fmt.Println("Part A:")
	partA := 0
	for _, lock := range locks {
		for _, key := range keys {
			if willFit(key, lock) {
				partA++
			}
		}
	}
	fmt.Printf("Number of unique key/lock pairs that fit: %d\n\n", partA)
}

func parseKeyLock(input []string) (keyLock []int, isKey bool) {
	isKey = input[0][0] == '.'
	if isKey {
		for i := 0; i < len(input)/2; i++ {
			input[i], input[len(input)-1-i] = input[len(input)-1-i], input[i]
		}
	}
	for i := 0; i < len(input[0]); i++ {
		for j := 1; j < len(input); j++ {
			if input[j][i] == '.' {
				keyLock = append(keyLock, j-1)
				break
			}
		}
	}
	return
}

func willFit(key, lock []int) bool {
	for i := 0; i < len(key); i++ {
		if key[i] > (5 - lock[i]) {
			return false
		}
	}
	return true
}
