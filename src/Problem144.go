package main

import (
	"fmt"
	"math"
)

func getL(m float64, b float64) (x float64, y float64) {
	x = (-b*m - 2*math.Sqrt(100-(b*b)+(25*m*m))) / (4 + m*m)
	y = m*x + b
	return
}

func getR(m float64, b float64) (x float64, y float64) {
	x = (-b*m + 2*math.Sqrt(100-(b*b)+(25*m*m))) / (4 + m*m)
	y = m*x + b
	return
}

func normal(x float64, y float64) float64 {
	return y / (4 * x)
}

func reflect(n float64, m float64) float64 {

	return (2*n - m*(1-(n*n))) / (1 + (2 * m * n) - (n * n))
}

func findB(x float64, y float64, m float64) float64 {
	return (-m * x) + y
}

func main() {

	m := -(197.0 / 14.0)
	b := 10.1

	x, y := getR(m, b)
	//fmt.Println("{", x, ",", y, "}, ")

	count := 0
	for math.Abs(x) > .01 || y < 0 {

		m = reflect(normal(x, y), m)

		b = findB(x, y, m)

		l, _ := getL(m, b)
		r, _ := getR(m, b)

		if math.Abs(l-x) > math.Abs(r-x) {
			x, y = getL(m, b)
		} else {
			x, y = getR(m, b)
		}

		//fmt.Println("{", x, ",", y, "}, ")

		count++

	}

	fmt.Println(count)
}
