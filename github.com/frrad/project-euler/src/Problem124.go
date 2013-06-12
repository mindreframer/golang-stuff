package main

import (
	"euler"
	"fmt"
	"time"
)

const (
	max    = 100000
	target = 10000
)

func main() {
	starttime := time.Now()

	radTable := make(map[int64][]int64)

	for i := int64(2); i <= max; i++ {
		factor := euler.Factor(i)
		factor = euler.RemoveDuplicates(factor)

		rad := int64(1)

		for _, fac := range factor {
			rad *= fac
		}

		radTable[rad] = append(radTable[rad], i)

	}

	current := 1
	for i := int64(0); i <= max; i++ {
		if answer, ok := radTable[i]; ok {
			//	fmt.Println(current)

			if current < target && current+len(answer) >= target {
				pos := target - current - 1
				fmt.Println(answer[pos])
			}
			//fmt.Println(i, answer)
			current += len(answer)
		}
	}

	fmt.Println("Elapsed time:", time.Since(starttime))
}
