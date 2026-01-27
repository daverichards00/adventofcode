package main

import (
	"fmt"
	"github.com/daverichards00/adventofcode/internal/convert"
	"github.com/daverichards00/adventofcode/internal/file"
	"github.com/daverichards00/adventofcode/internal/maths"
	"slices"
)

const fftPhaseFreq = 100

var (
	basePattern = []int{0, 1, 0, -1}
)

func main() {
	fmt.Println("Day 16")

	input := file.Load("cmd/2019/day16/input.txt")
	signal := convert.StrToDigits(input[0])

	fmt.Println("Part A:")
	outputA := fftv1(signal)
	fmt.Printf("8 digit message: %s\n\n", convert.DigitsToStr(outputA))

	fmt.Println("Part B:")
	fullSignal := slices.Repeat(signal, 10000)
	outputB := fftv2(fullSignal)
	fmt.Printf("8 digit message: %s\n\n", convert.DigitsToStr(outputB))

}

func getPattern(length int) [][]int {
	lenBase := len(basePattern)
	p := make([][]int, length)
	for i := range p {
		p[i] = make([]int, length)
		for j := range p[i] {
			p[i][j] = basePattern[((j+1)/(i+1))%lenBase]
		}
	}
	return p
}

func fftv1(signal []int) []int {
	pattern := getPattern(len(signal))

	for i := 0; i < fftPhaseFreq; i++ {
		next := make([]int, len(signal))
		for j := range next {
			for k := range signal {
				next[j] += signal[k] * pattern[j][k]
			}
			next[j] = maths.Abs(next[j]) % 10
		}
		signal = next
	}
	return signal[0:8]
}

func fftv2(signal []int) []int {
	msgIdx := convert.DigitsToInt(signal[0:7])
	if msgIdx < len(signal)/2 {
		panic("unsupported message index")
	}

	output := signal[msgIdx:]
	slices.Reverse(output)

	for p := 0; p < fftPhaseFreq; p++ {
		cumsum := make([]int, len(output))
		cumsum[0] = output[0]
		for i := 1; i < len(output); i++ {
			cumsum[i] = (cumsum[i-1] + output[i]) % 10
		}
		output = cumsum
	}

	msg := output[len(output)-8:]
	slices.Reverse(msg)
	return msg
}
