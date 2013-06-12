package genetic

import (
	rnd "github.com/handcraftsman/Random"
	"math/rand"
	"runtime"
	s "sort"
	"time"
)

func createRandomNumberGenerator() randomSource {
	procs := runtime.GOMAXPROCS(-1)
	if procs > 1 {
		return rnd.NewRandom()
	}
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func insertionSort(items []*sequenceInfo, compare func(*sequenceInfo, *sequenceInfo) bool, index int) {
	if index < 1 || index > len(items) {
		return
	}

	if index > 0 {
		location := s.Search(index, func(i int) bool { return compare(items[index], items[i]) })
		temp := items[index]
		for i := index; i > location; i-- {
			items[i] = items[i-1]
		}
		items[location] = temp
	}
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

func reverseArray(a []string) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}

func sort(a, b int) (int, int) {
	if a < b {
		return a, b
	}
	return b, a
}
