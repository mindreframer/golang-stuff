package main

import (
	"./euler"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

func validateSub(number, word string) bool {
	if len(number) != len(word) {
		return false
	}
	if number[0:1] == "0" {
		return false
	}

	digit := make(map[string]string)
	letter := make(map[string]string)

	for i := 0; i < len(word); i++ {

		if answer, ok := digit[number[i:i+1]]; ok && answer != word[i:i+1] {
			return false
		}
		if answer, ok := letter[word[i:i+1]]; ok && answer != number[i:i+1] {
			return false
		}
		letter[word[i:i+1]] = number[i : i+1]
		digit[number[i:i+1]] = word[i : i+1]
	}

	return true
}

func main() {
	starttime := time.Now()

	data := euler.Import("problemdata/words.txt")
	line := data[0]
	words := strings.Split(line, ",")

	anagrams := make(map[string][]string)

	for _, word := range words {
		stripped := word[1 : len(word)-1]
		sorted := euler.BubbleSort(stripped)

		anagrams[sorted] = append(anagrams[sorted], stripped)
	}

	max := 0

	for _, set := range anagrams {
		if len(set) > 1 {
			//We're ignoring triples, etc (data contains only 1 and it's short)
			word1 := set[0]
			word2 := set[1]
			length := len(word1)

			permutation := euler.UnPermuteStrings(word1, word2)

			lower := int(math.Sqrt(float64(euler.IntExp(10, int64(length)-1))))
			upper := int(math.Sqrt(float64(euler.IntExp(10, int64(length)))))

			for i := lower; i < upper; i++ {
				square1 := i * i
				s1string := strconv.Itoa(square1)
				s2string := euler.PermuteString(permutation, s1string)
				square2, _ := strconv.Atoi(s2string)
				if validateSub(s1string, word1) && validateSub(s2string, word2) && euler.IsSquare(int64(square2)) {

					if square1 > max {
						max = square1
					}
					if square2 > max {
						max = square2
					}
				}
			}
		}
	}

	fmt.Println(max)

	fmt.Println("Elapsed time:", time.Since(starttime))

}
