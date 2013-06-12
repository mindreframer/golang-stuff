package main

import (
	"./euler"
	"fmt"
	"time"
)

const (
	target = 100000000
)

func main() {

	starttime := time.Now()

	euler.PrimeCache(target)

	counter := int64(0)

	for i := int64(1); euler.Prime(i)*euler.Prime(i) < target; i++ {

		counter += euler.PrimePi(target/euler.Prime(i)) - i + 1

	}

	fmt.Println(counter)

	fmt.Println("Elapsed time:", time.Since(starttime))

}
