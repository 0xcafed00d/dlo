package main

import (
	"math/rand"
	"sync"
	"time"
)

type AtomicRand struct {
	sync.Mutex
	r *rand.Rand
}

func MakeAtomicRand() *AtomicRand {
	return &AtomicRand{r: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

func (ar *AtomicRand) Int63n(n int64) int64 {
	ar.Lock()
	defer ar.Unlock()

	return ar.r.Int63n(n)

}
