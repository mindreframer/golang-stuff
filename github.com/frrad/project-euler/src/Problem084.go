package main

import (
	"fmt"
	"math/rand"
	"time"
)

func roll() (int, bool) {
	r1, r2 := rand.Int()%4+1, rand.Int()%4+1
	return r1 + r2, r1 == r2
}

func main() {
	starttime := time.Now()
	board := make([]int, 40)
	place := 0
	total := 0

	var jailtimer int

	for i := 0; i < 10000000; i++ {

		board[place]++
		total++

		role, doubles := roll()
		if doubles {
			jailtimer++
		} else {
			jailtimer = 0
		}

		place = (place + role) % 40

		if jailtimer == 3 {
			place = 10
			jailtimer = 0
		}

		if place == 30 {
			//G2J i=30
			place = 10
		} else if place == 2 || place == 17 || place == 33 {
			//CC i=2,17,33
			dice := rand.Int() % 16
			if dice == 1 {
				place = 0
			} else if dice == 2 {
				place = 10
			}

		} else if place == 7 || place == 22 || place == 36 {

			//CH i=7,22,36
			switch dice := rand.Int() % 16; dice {
			case 1:
				place = 0
			case 2:
				place = 10
			case 3:
				place = 11
			case 4:
				place = 24
			case 5:
				place = 39
			case 6:
				place = 5
			case 7:
				if place == 7 {
					place = 15
				}
				if place == 22 {
					place = 25
				}
				if place == 36 {
					place = 5
				}
			case 8:
				if place == 7 {
					place = 15
				}
				if place == 22 {
					place = 25
				}
				if place == 36 {
					place = 5
				}
			case 9:
				if place == 7 {
					place = 12
				}
				if place == 22 {
					place = 28
				}
				if place == 36 {
					place = 12
				}
			case 10:
				place = place - 3
			}

		}

	}

	for j := 0; j < 3; j++ {

		max := 0
		maxi := 0

		for i, spot := range board {
			if spot > max {
				max = spot
				maxi = i
			}

		}

		board[maxi] = 0

		fmt.Println(maxi, ":", 100*float64(max)/float64(total))
	}

	fmt.Println("Elapsed time:", time.Since(starttime))
}
