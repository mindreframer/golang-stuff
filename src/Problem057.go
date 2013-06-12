package main

import (
	"./euler"
	"fmt"
)

func main() {

	num := "2"

	den := "1"

	total := 0
	counter := 2

	for i := 0; i < 1005; i++ {

		num, den = euler.StringFastFracAdd("2", "1", den, num)
		if len(euler.StringSum(num, den)) > len(num) {
			total++
			fmt.Println(euler.StringSum(num, den), num, ": (", total, "/", counter, ")")

		}
		counter++

	}

}
