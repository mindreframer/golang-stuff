package main

import (
	"./eulerlib"
	"fmt"
	"strconv"
)

func intSubstr(number int64, place int) int64 {
	str := strconv.FormatInt(number, 10)
	stri := ""

	for i := 0; i < 3; i++ {
		stri += string(str[i+place])
	}

	swapped, _ := strconv.ParseInt(stri, 10, 64)

	return swapped

}

func main() {

	total := int64(0)
	for i := int64(1000000000); i < 9999999999; i++ {
		if intSubstr(i, 1)%2 == 0 && intSubstr(i, 2)%3 == 0 &&
			intSubstr(i, 3)%5 == 0 && intSubstr(i, 4)%7 == 0 && intSubstr(i, 5)%11 == 0 &&
			intSubstr(i, 6)%13 == 0 && intSubstr(i, 7)%17 == 0 && eulerlib.IsPandigital(i) {
			fmt.Println(i, total)
			total += i
		}

		if i%1000000 == 0 {
			fmt.Println(i)
		}
	}
	fmt.Println(total)

}
