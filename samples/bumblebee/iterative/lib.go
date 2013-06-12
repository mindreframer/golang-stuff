package main

import (
	"math/rand"
	"time"
)

func createRandomNumberGenerator() randSource {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

type randSource interface {
	Intn(int) int
}

func getUniqueInt(rand randSource, used map[int]bool, max int) (value int) {
	value = 0
	for {
		value = rand.Intn(max)
		if !used[value] {
			break
		}
	}
	return
}

func sort(a, b int) (int, int) {
	if a < b {
		return a, b
	}
	return b, a
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
