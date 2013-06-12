package main

import (
	"./euler"
	"fmt"
	"time"
)

const (
	hat = 50000000
)

func main() {
	starttime := time.Now()

	table := [hat]int{}

	pre := int64(0)

	for i := int64(1); pre < hat; i++ {
		pre = euler.Prime(i) * euler.Prime(i) * euler.Prime(i) * euler.Prime(i)
		pri := int64(0)

		for j := int64(1); pri < hat; j++ {
			pri = pre + euler.Prime(j)*euler.Prime(j)*euler.Prime(j)

			pro := int64(0)

			for k := int64(1); pro < hat; k++ {
				pro = pri + euler.Prime(k)*euler.Prime(k)

				if pro < hat {
					table[pro]++
				}
			}

		}

	}

	total := 0
	for i := 0; i < hat; i++ {

		if table[i] > 0 {
			total++
		}
	}

	fmt.Println(total)

	fmt.Println("Elapsed time:", time.Since(starttime))

}
