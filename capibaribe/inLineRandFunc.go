package capibaribe

import (
	"math/rand"
	"time"
)

func inLineRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}
