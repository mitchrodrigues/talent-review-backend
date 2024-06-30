package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomFloatBetween(t *testing.T) {
	min := 10.0
	max := 20.0
	const runs = 100

	for i := 0; i < runs; i++ {
		result := RandomFloatBetween(min, max)
		assert.GreaterOrEqual(t, result, min, "Result should be greater than or equal to min")
		assert.Less(t, result, max, "Result should be less than max")
	}
}
