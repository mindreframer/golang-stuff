package main

import (
	"fmt"
	"time"
)

func basisPoly(x float64, index int, points [][2]float64) float64 {
	answer := float64(1)
	for i := 0; i < index; i++ {
		answer *= (x - points[i][0]) / (points[index][0] - points[i][0])
	}
	for i := index + 1; i < len(points); i++ {
		answer *= (x - points[i][0]) / (points[index][0] - points[i][0])
	}
	return answer
}

func lagrange(x float64, points [][2]float64) float64 {
	answer := float64(0)
	for i := 0; i < len(points); i++ {
		answer += points[i][1] * basisPoly(x, i, points)
	}
	return answer
}

func f(n float64) float64 {
	return 1 - n + (n * n) - (n * n * n) + (n * n * n * n) - (n * n * n * n * n) + (n * n * n * n * n * n) - (n * n * n * n * n * n * n) + (n * n * n * n * n * n * n * n) - (n * n * n * n * n * n * n * n * n) + (n * n * n * n * n * n * n * n * n * n)
}

func main() {
	starttime := time.Now()

	pts := [][2]float64{[2]float64{1, f(1)}}

	answer := float64(0)

	for i := float64(0); i < 10; i++ {
		answer += lagrange(i+2, pts)
		pts = append(pts, [2]float64{i + 2, f(i + 2)})
	}

	fmt.Println(int64(answer))

	fmt.Println("Elapsed time:", time.Since(starttime))
}
