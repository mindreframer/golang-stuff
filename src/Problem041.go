package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func sortInt(input int) int {

	swapped, _ := strconv.Atoi(bubbleSort(strconv.Itoa(input)))
	return swapped

}

func bubbleSort(word string) string {
	wordtable := strings.Split(word, "")
	for j := 0; j < len(word); j++ {

		for i := 0; i < len(word)-1; i++ {
			if wordtable[i] < wordtable[i+1] {
				temp := wordtable[i]
				wordtable[i] = wordtable[i+1]
				wordtable[i+1] = temp
			}
		}
	}
	return strings.Join(wordtable, "")
}

func isPrime(n int) bool {
	if n == 1 {
		return false
	}
	lim := int(math.Sqrt(float64(n))) + 1
	for i := 2; i < lim; i++ {
		if n%i == 0 {
			return false
		}
	}

	return true
}

func isPandigital(n int) bool {

	height := 1 + int(math.Log10(float64(n)))

	output := 0

	for i := 1; i < height+1; i++ {
		current := 1
		for j := 1; j < i; j++ {
			current *= 10
		}
		output += (current * i)
	}

	return output == sortInt(n)
}

func main() {
	for i := 0; i < 999999999; i++ {
		if isPandigital(i) && isPrime(i) {
			fmt.Println(i)
		}
	}

}
