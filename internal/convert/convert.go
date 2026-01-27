package convert

import "strconv"

func DigitsToInt(digits []int) int {
	if len(digits) == 0 {
		return 0
	}
	i := digits[0]
	for _, d := range digits[1:] {
		i = (i * 10) + d
	}
	return i
}

func DigitsToStr(digits []int) string {
	str := make([]byte, len(digits))
	for i := range digits {
		str[i] = '0' + byte(digits[i])
	}
	return string(str)
}

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
