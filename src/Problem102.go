package main

import (
	"euler"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type line struct {
	m float64
	b float64
}

type point struct {
	x float64
	y float64
}

func onLeft(a point, f line) bool {
	return a.y < (f.m*a.x)+f.b
}

func makeLine(a, b point) line {
	M := (a.y - b.y) / (a.x - b.x)
	B := a.y - (a.x * M)
	return line{M, B}
}

func main() {

	starttime := time.Now()

	origin := point{0, 0}
	data := euler.Import("../problemdata/triangles.txt")
	total := 0


	for _, line := range data {

		split := strings.Split(line, ",")
		numbers := make([]float64, 6)
		for i, number := range split {
			temp, _ := strconv.Atoi(number)
			numbers[i] = float64(temp)
		}

		point1 := point{numbers[0], numbers[1]}
		point2 := point{numbers[2], numbers[3]}
		point3 := point{numbers[4], numbers[5]}

		f1 := makeLine(point1, point2)
		f2 := makeLine(point2, point3)
		f3 := makeLine(point3, point1)

		inside := onLeft(origin, f1) == onLeft(point3, f1)
		inside = inside && onLeft(origin, f2) == onLeft(point1, f2)
		inside = inside && onLeft(origin, f3) == onLeft(point2, f3)

		if inside {
			total++
		}
	}

	fmt.Println(total)

	fmt.Println("Elapsed time:", time.Since(starttime))
}
