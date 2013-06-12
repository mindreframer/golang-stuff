package main

import (
	"euler"
	"fmt"
	"time"
)

const (
	lid    = 1000000000000000000
	target = 30
)

func main() {
	starttime := time.Now()

	total := 0
	list := make([]int64, 0)

	for j := int64(2); j < 100; j++ {
		for i := int64(2); euler.IntExp(j, i) < lid && euler.IntExp(j, i) > 0; i++ {
			if int64(euler.DigitSum((euler.IntExp(j, i)))) == j {

				total++
				list = append(list, euler.IntExp(j, i))
			}
		}
	}

	list = euler.SortLInts(list)

	euler.ReverseLInts(list)

	fmt.Println(list[target-1])

	fmt.Println("Elapsed time:", time.Since(starttime))

}
