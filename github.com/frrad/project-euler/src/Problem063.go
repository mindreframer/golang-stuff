package main

import (
	"./euler"
	"fmt"
)

func main() {
	total := 0
	for i := int64(1); i < 10; i++ {
		for j := int64(1); j < 26; j++ {

			temp := euler.IntExp(i, j)
			if int64(euler.NumberDigits(temp)) == j {
				fmt.Println(i, j, temp)
				total++
			}
		}
	}

	fmt.Println("Total:", total)
	//Answer is 49: start to int overflow at the end.

}
