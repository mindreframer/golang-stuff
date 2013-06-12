package main

import (
	"./euler"
	"fmt"
	"time"
)

func main() {
	starttime := time.Now()

	sum := int64(0)

	for i := int64(2); i < 1000001; i++ {
		sum += euler.Totient(i)
	}

	fmt.Println(sum)

	fmt.Println("Elapsed time:", time.Since(starttime))

}
