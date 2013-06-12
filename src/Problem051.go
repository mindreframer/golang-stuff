package main

import (
	"./euler"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func main() {
	starttime := time.Now()

	i := int64(0)
	score := 0
	for score < 8 {
		pattern := strconv.FormatInt(i, 11)
		if strings.Count(pattern, "a") != 0 {

			score = 0
			for i := 0; i < 10; i++ {
				I := strconv.Itoa(i)
				n, _ := strconv.Atoi(strings.Replace(pattern, "a", I, -1))
				if euler.IsPrime(int64(n)) {
					score++
				}
			}

			if score > 7 {
				fmt.Println(pattern, score)
			}

		}
		i++
	}

	fmt.Println("Elapsed time:", time.Since(starttime))

}
