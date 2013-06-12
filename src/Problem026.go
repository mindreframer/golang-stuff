package main

import "fmt"

const level = 10


func chomp(numerator int, denominator int, height int) int {

	for i := 1; i <= height; i++ {

		numerator *= 10
		numerator = numerator % denominator

	}

	return numerator
}

func main() {

	record := 0

	for den := 2; den < 1000; den++ {

		newmerator := chomp(1, den, level)

		answer := 0
		for j := 1; newmerator != chomp(newmerator, den, j); j++ {//this is pretty wasteful
			answer = j + 1
		}
		if answer > record {
			record = answer
			fmt.Println(den, ",", answer)
		}

	}

}

