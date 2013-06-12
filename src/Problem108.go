package main

import (
	"fmt"
	"time"
)

func main() {
	starttime := time.Now()

	target := int64(4)

	counter := int64(0)
	record := int64(0)

	for target = 4; counter < 1000; target++ {

		counter = 0

		for den := target + 1; den <= 2*target; den++ {

			if (den*target)%(den-target) == 0 {
				counter++
			}

		}

		if counter > record {
			record = counter
			fmt.Println(target, record)
		}

	}

	fmt.Println(target)

	fmt.Println("Elapsed time:", time.Since(starttime))
}
