package math

import "math"

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Pow(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}
