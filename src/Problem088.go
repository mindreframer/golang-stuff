package main

import (
	"euler"
	"fmt"
	"time"
)

func sum(list []int) (total int) {
	for _, item := range list {
		total += item
	}
	return
}

func prod(list []int) (prod int) {
	prod = 1
	for _, item := range list {
		prod *= item
	}
	return
}

const (
	infinity = 99999999
)

func search(length int, c chan int) {

	limit := 3 * length

	test := make([]int, length)

	inner := func() {
		//fmt.Println(test, limit)
		if sum(test) == prod(test) {
			//fmt.Println(test)
			limit = sum(test)
		}
	}

	var level func(int)

	level = func(lvl int) {

		start := 1

		if lvl > 1 {
			start = test[lvl-2]
		}

		leftprod := int64(1)

		var rightprod int64
		if start >= 2 && 1+length-lvl >= 13 {
			rightprod = infinity
		} else {
			rightprod = euler.IntExp(int64(start), int64(1+length-lvl))
		}

		for test[lvl-1] = start; sum(test[0:lvl])+((length-lvl)*test[lvl-1]) < limit && leftprod*rightprod < int64(limit); test[lvl-1]++ {

			if lvl >= length {
				inner()
			} else {
				level(lvl + 1)
			}

			if lvl > 0 {
				leftprod = int64(prod(test[0 : lvl-1]))
			} else {
				leftprod = 1
			}

			if test[lvl-1] >= 2 && 1+length-lvl >= 13 {
				rightprod = infinity
			} else {
				rightprod = euler.IntExp(int64(test[lvl-1]), int64(1+length-lvl))
			}

			//fmt.Println(test[lvl-1], "to the", length-lvl+1, "is", rightprod)
		}
	}

	level(1)

	c <- limit
}

func main() {
	starttime := time.Now()

	recordTable := make(map[int]bool)

	bottom := 2
	top := 12000

	c := make(chan int)

	for length := bottom; length <= top; length++ {

		go search(length, c)
		fmt.Println("launched", length)

	}

	for length := bottom; length <= top; length++ {

		answer := <-c
		fmt.Println(length)
		fmt.Println(answer)

		recordTable[answer] = true
	}

	total := 0
	for i := range recordTable {
		total += i
	}

	fmt.Println(total)

	fmt.Println("Elapsed time:", time.Since(starttime))
}
