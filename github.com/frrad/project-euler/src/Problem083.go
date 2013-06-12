package main

import (
	"euler"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func main() {
	starttime := time.Now()

	size := 80
	data := euler.Import("../problemdata/matrix.txt")
	matrix := make([][]int, len(data))

	for j, line := range data {
		words := strings.Split(line, ",")
		row := make([]int, len(words))
		for i, word := range words {
			row[i], _ = strconv.Atoi(word)
		}
		matrix[j] = row
	}

	dist := make([][]int, size)
	for i := 0; i < size; i++ {
		dist[i] = make([]int, size)
	}

	dist[0][0] = matrix[0][0]

	for i := 1; i < size; i++ {
		dist[0][i] = matrix[0][i] + dist[0][i-1]
	}

	for j := 1; j < size; j++ {
		for i := 0; i < size; i++ {
			dist[j][i] = matrix[j][i] + dist[j-1][i]
		}

	}

	currentSum := 1
	lastSum := 0

	for currentSum != lastSum {

		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {

				if i != 0 {
					if matrix[i][j]+dist[i-1][j] < dist[i][j] {
						dist[i][j] = matrix[i][j] + dist[i-1][j]
					}
				}

				if j != 0 {
					if matrix[i][j]+dist[i][j-1] < dist[i][j] {
						dist[i][j] = matrix[i][j] + dist[i][j-1]
					}
				}

				if i != size-1 {
					if matrix[i][j]+dist[i+1][j] < dist[i][j] {
						dist[i][j] = matrix[i][j] + dist[i+1][j]
					}
				}

				if j != size-1 {
					if matrix[i][j]+dist[i][j+1] < dist[i][j] {
						dist[i][j] = matrix[i][j] + dist[i][j+1]
					}
				}

			}
		}

		lastSum = currentSum
		currentSum = 0
		for j := 1; j < size; j++ {
			for i := 0; i < size; i++ {
				currentSum += dist[i][j]
			}

		}

	}

	fmt.Println(dist[size-1][size-1])

	fmt.Println("Elapsed time:", time.Since(starttime))
}
