package main

import (
	"fmt"
	"strconv"
	"time"
)

func fitsMask(n int64) bool {
	w := strconv.FormatInt(n, 10)

	if len(w) != 19 {
		return false
	}

	for i := 1; i < 10; i++ {
		if w[2*i-2:2*i-1] != strconv.Itoa(i) {
			return false
		}
	}

	if w[18:] != "0" {
		return false
	}

	return true
}

func main() {
	starttime := time.Now()

	solution := int64(0)

	for i := int64(1000000000); !fitsMask(i * i); i += 10 {

		solution = i + 10

	}
	fmt.Println(solution, solution*solution)

	fmt.Println("Elapsed time:", time.Since(starttime))

}
