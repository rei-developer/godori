package math

import (
	"math"
	"math/rand"
	"time"
)

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

func Rand(val int) int {
	rs := rand.NewSource(time.Now().UnixNano())
	r := rand.New(rs)
	return r.Intn(val)
}
