package main

import (
	"./euler"
	"fmt"
	"strconv"
)

func ctdFrac(list []int) (num string, den string) {
	num = strconv.Itoa(list[len(list)-1])
	den = "1"

	for i := len(list) - 2; i >= 0; i-- {

		num, den = euler.StringFastFracAdd(strconv.Itoa(list[i]), "1", den, num)

		fmt.Println(num, den)
	}

	return
}

func eList(n int) []int {
	answer := make([]int, n)

	answer[0] = 2

	for i := 1; i < n; i++ {
		answer[i] = 1
	}

	for i := 0; 3*i+2 < n; i++ {
		answer[3*i+2] = 2 * (i + 1)

	}

	return answer
}

func main() {

	fmt.Println(euler.StringSum("0", "2345654322"))

	numerator, _ := ctdFrac(eList(100))
	fmt.Println(euler.StringDigitSum(numerator))

}
