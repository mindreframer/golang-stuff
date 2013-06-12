package main

import (
	"./euler"
	"fmt"
)

func evalNumeral(numeral string) int {
	if len(numeral) == 0 {
		return 0
	}
	if len(numeral) == 1 {
		if numeral == "I" {
			return 1
		}
		if numeral == "V" {
			return 5
		}
		if numeral == "X" {
			return 10
		}
		if numeral == "L" {
			return 50
		}
		if numeral == "C" {
			return 100
		}
		if numeral == "D" {
			return 500
		}
		if numeral == "M" {
			return 1000
		}
	}

	a := numeral[len(numeral)-2 : len(numeral)-1]
	b := numeral[len(numeral)-1:]
	rest := numeral[:len(numeral)-2]

	A := evalNumeral(a)
	B := evalNumeral(b)

	if A < B {
		return evalNumeral(rest) + B - A
	}

	return B + evalNumeral(rest+a)
}

func best(n int) int {

	char1 := []int{1, 5, 10, 50, 100, 500, 1000}
	char2 := []int{4, 9, 40, 90, 400, 900}

	if n == 0 {
		return 0
	}

	a, b := 1, 0

	for i := 0; i < len(char1); i++ {
		if char1[i] <= n {
			a = char1[i]
		}
	}
	for i := 0; i < len(char2); i++ {
		if char2[i] <= n {
			b = char2[i]
		}
	}

	if b != 0 && b > a {
		return 2 + best(n-b)
	}

	return 1 + best(n-a)
}

func main() {

	data := euler.Import("problemdata/roman.txt")

	blerg := 0
	for _, numeral := range data {

		blerg += len(numeral) - best(evalNumeral(numeral))
	}
	fmt.Println(blerg)

}
