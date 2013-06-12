package main

import (
	"./euler"
	"fmt"
)

func main() {

	record := 0.

	for i := int64(2); i < 1000000; i++ {
		ratio := float64(i) / float64(euler.Totient(i))
		if ratio > record {
			fmt.Println(i, euler.Totient(i), ratio)
			record = ratio
		}
	}

}
