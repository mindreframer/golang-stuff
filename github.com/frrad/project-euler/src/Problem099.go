package main

import (
	"./euler"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

//returns true if a^b < c^d
func biggerThan(a, b, c, d float64) bool {
	return math.Log(a)*b > math.Log(c)*d

}

func main() {
	starttime := time.Now()

	data := euler.Import("problemdata/base_exp.txt")

	maxa := 1.
	maxb := 1.

	for i, line := range data {
		linesplit := strings.Split(line, ",")

		a, _ := strconv.Atoi(linesplit[0])
		b, _ := strconv.Atoi(linesplit[1])

		if biggerThan(float64(a), float64(b), maxa, maxb) {
			fmt.Println(line, i+1)
			maxa = float64(a)
			maxb = float64(b)
		}

	}

	fmt.Println("Elapsed time:", time.Since(starttime))

}
