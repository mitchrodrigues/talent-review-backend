package helpers

import (
	"math/rand"
)

var (
	localRand = rand.New(rand.NewSource(99))
)

func RandomNumberBetween(min, max int) int {
	if min == max {
		return max
	}

	if min > max {
		min, max = max, min // Swap if min is greater than max
	}

	return localRand.Intn(max-min+1) + min
}

func RandomFloatBetween(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
