package main

import (
	"fmt"
	"time"
)

func prob(won, drawn, rounds int) float64 {

	if drawn == rounds {
		if 2*won > drawn {
			return 1
		} else {
			return 0
		}
	}

	return prob(won+1, drawn+1, rounds)*1/float64(drawn+2) + prob(won, drawn+1, rounds)*float64(drawn+1)/float64(drawn+2)

}

func main() {
	starttime := time.Now()

	fmt.Println(int(1 / prob(0, 0, 15)))

	fmt.Println("Elapsed time:", time.Since(starttime))
}
