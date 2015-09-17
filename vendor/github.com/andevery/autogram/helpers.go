package autogram

import (
	"log"
	"math/rand"
)

func randomIndexes(length, count int) []int {
	log.Println(length, count)
	if count > length {
		count = length
	}
	return rand.Perm(length)[:count]
}
