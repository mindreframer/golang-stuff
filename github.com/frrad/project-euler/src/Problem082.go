package main

import (
	"./euler"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const infinity = 999999

func sign(n int) int {
	if n < 0 {
		return -1
	}
	return 1
}

func main() {
	starttime := time.Now()

	data := euler.Import("problemdata/matrix.txt")

	matrix := make([][]int, len(data))

	for j, line := range data {
		words := strings.Split(line, ",")
		row := make([]int, len(words))
		for i, word := range words {
			row[i], _ = strconv.Atoi(word)
		}
		matrix[j] = row
	}

	sumtrix := make([][]int, len(matrix))
	for i := range sumtrix {
		sumtrix[i] = make([]int, len(matrix[0]))
		sumtrix[i][0] = matrix[i][0]
	}

	for currentColumn := 1; currentColumn < len(matrix[0]); currentColumn++ {

		for i := 0; i < len(matrix); i++ {
			winner := infinity
			for j := 0; j < len(matrix); j++ {
				//now we compute the cost of getting to (i, column) from (j, column -1)
				cost := matrix[i][currentColumn]
				cost += sumtrix[j][currentColumn-1]
				for k := i; sign(j-i)*(j-k) > 0; k += sign(j - i) {
					cost += matrix[k][currentColumn-1]
				}

				if cost < winner {
					winner = cost
				}

			}

			sumtrix[i][currentColumn] = winner

		}
	}

	best := infinity

	for i := 0; i < len(sumtrix); i++ {
		if sumtrix[i][len(sumtrix[i])-1] < best {
			best = sumtrix[i][len(sumtrix[i])-1]
		}
	}

	fmt.Println(best)

	fmt.Println("Elapsed time:", time.Since(starttime))

}
