package main

import (
	"fmt"
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

func main() {

	for i := 0; ; i++ {
		bench := sortInt(i)
		if sortInt(2*i) == bench &&
			sortInt(3*i) == bench &&
			sortInt(4*i) == bench &&
			sortInt(5*i) == bench &&
			sortInt(6*i) == bench {
			fmt.Println(i)
		}
	}

}
