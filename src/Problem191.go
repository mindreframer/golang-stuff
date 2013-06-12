package main

import (
	"fmt"
	"strconv"
	"time"
)

var memo map[string]int

func numb(length int, leading2 string, hasbeenlate bool) int {

	if length == 2 {
		if leading2 == "AA" {
			return 1
		}
		if leading2 == "AL" {
			if hasbeenlate {
				return 0
			}
			return 1
		}
		if leading2 == "AO" {
			return 1
		}
		if leading2 == "LA" {
			if hasbeenlate {
				return 0
			}
			return 1
		}
		if leading2 == "LL" {
			return 0
		}
		if leading2 == "LO" {
			if hasbeenlate {
				return 0
			}
			return 1
		}
		if leading2 == "OA" {
			return 1

		}
		if leading2 == "OL" {
			if hasbeenlate {
				return 0
			}
			return 1
		}
		if leading2 == "OO" {
			return 1
		}
	}

	if answer, ok := memo[keygen(length, leading2, hasbeenlate)]; ok {
		return answer
	}

	a := leading2[0:1]
	b := leading2[1:2]

	if (a == "L" || b == "L") && hasbeenlate {
		return 0
	}

	answer := 0

	nextlate := hasbeenlate
	if a == "L" {
		nextlate = true
	}

	if a == b && a == "A" {

	} else {
		answer += numb(length-1, b+"A", nextlate)
	}

	answer += numb(length-1, b+"L", nextlate)
	answer += numb(length-1, b+"O", nextlate)

	memo[keygen(length, leading2, hasbeenlate)] = answer

	return answer

}

func keygen(n int, s string, b bool) string {
	answer := ""
	answer += strconv.Itoa(n)
	answer += s
	if b {
		answer += "T"
	} else {
		answer += "F"
	}
	return answer
}

func main() {
	starttime := time.Now()

	memo = make(map[string]int)
	total := 30

	fmt.Println(numb(total, "OO", false) + numb(total, "OA", false) + numb(total, "OL", false) + numb(total, "AO", false) + numb(total, "AA", false) + numb(total, "AL", false) + numb(total, "LO", false) + numb(total, "LA", false))

	fmt.Println("Elapsed time:", time.Since(starttime))
}
