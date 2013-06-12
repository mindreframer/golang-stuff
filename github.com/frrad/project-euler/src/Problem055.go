package main

import (
	"./eulerlib"
	"fmt"
	"strconv"
)

func isl(number string) bool {
	current := number
	for i := 0; i < 50; i++ {
		current = eulerlib.StringSum(current, eulerlib.StringReverse(current))
		if eulerlib.IsStringPalindrome(current) {
			return false
		}
	}

	return true
}

func main() {

	total := 0
	for i := 0; i < 10000; i++ {
		if isl(strconv.Itoa(i)) {
			fmt.Println(i)
			total++
		}
	}

	fmt.Println("Total:", total)
}
