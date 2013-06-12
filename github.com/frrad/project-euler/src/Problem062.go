package main

import (
	"./euler"
	"fmt"
	"time"
)

func main() {
	starttime := time.Now()

	var cubes = map[int64]int{}

	x := int64(1)
	i := int64(1)

	for ; cubes[x] < 5; i++ {
		x = euler.SortInt(i * i * i)

		cubes[x]++

	}

	x = int64(1)
	i = int64(1)

	for ; cubes[x] != 5; i++ {
		x = euler.SortInt(i * i * i)

	}

	i = i - 1

	fmt.Println(i * i * i)

	fmt.Println("Elapsed time:", time.Since(starttime))

}
