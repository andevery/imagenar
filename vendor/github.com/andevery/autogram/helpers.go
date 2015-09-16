package autogram

import (
	"math/rand"
)

func randomIndexes(length, count int) []int {
	return rand.Perm(length)[:count]
}
