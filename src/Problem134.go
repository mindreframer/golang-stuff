package main

import (
	"./euler"
	"fmt"
	"time"
)

const top = 1000000

func pair(p1, p2 int64, c chan int64) {
	mod := int64(10)

	for ; mod < p1; mod *= 10 {
	}

	for try := int64(0); ; try += p2 {
		if try%mod == p1 {
			c <- try
			return
		}

	}

}

func main() {
	starttime := time.Now()

	euler.PrimeCache(top)

	c := make(chan int64)

	go func() {
		for i := int64(3); euler.Prime(i) < top; i++ {
			go pair(euler.Prime(i), euler.Prime(i+1), c)
		}
	}()

	total := int64(0)
	for i := int64(0); i < euler.PrimePi(top)-2; i++ {
		total += <-c
	}
	fmt.Println(total)

	fmt.Println("Elapsed time:", time.Since(starttime))

}
