package convert

import "strconv"

func StrToDigits(s string) []int {
	digits := make([]int, len(s))
	for i := range s {
		digits[i] = StrToInt(string(s[i]))
	}
	return digits
}

func StrToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
