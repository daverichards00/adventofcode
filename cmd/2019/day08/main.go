package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/internal/convert"
	"github.com/daverichards00/adventofcode/internal/file"
)

func main() {
	fmt.Println("Day 08")
	input := file.Load("cmd/2019/day08/input.txt")

	w, h := 25, 6

	fmt.Println("Part A:")
	layers := decodeLayers(input[0], w, h)
	fewest0Idx, fewest0Freq := 0, pixelFreq(layers[0], 0)
	for i := 1; i < len(layers); i++ {
		if f := pixelFreq(layers[i], 0); f < fewest0Freq {
			fewest0Idx, fewest0Freq = i, f
		}
	}
	fmt.Printf("Sum of 1 digits multiplied by sumer of 2 digits: %d\n\n", pixelFreq(layers[fewest0Idx], 1)*pixelFreq(layers[fewest0Idx], 2))

	fmt.Println("Part B:")
	img := renderImage(layers)
	printImage(img)
}

func decodeLayers(raw string, w, h int) [][][]int {
	l := make([][][]int, len(raw)/(w*h))
	p := 0
	for i := range l {
		l[i] = make([][]int, h)
		for j := range l[i] {
			l[i][j] = make([]int, w)
			for k := range l[i][j] {
				l[i][j][k] = convert.StrToInt(string(raw[p]))
				p++
			}
		}
	}
	return l
}

func renderImage(layers [][][]int) [][]int {
	img := make([][]int, len(layers[0]))
	for i := range img {
		img[i] = make([]int, len(layers[0][0]))
		for j := range img[i] {
			for _, layer := range layers {
				if layer[i][j] != 2 {
					img[i][j] = layer[i][j]
					break
				}
			}
		}
	}
	return img
}

func printImage(img [][]int) {
	for i := range img {
		for j := range img[i] {
			switch img[i][j] {
			case 1:
				fmt.Print("#")
			default:
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func pixelFreq(layer [][]int, p int) int {
	f := 0
	for _, row := range layer {
		for _, col := range row {
			if col == p {
				f++
			}
		}
	}
	return f
}
