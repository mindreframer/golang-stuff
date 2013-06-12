package main

import (
	"fmt"
	"time"
)

func works(a, b map[int]bool) bool {
	for i := 1; i <= 9; i++ {
		d1 := (i * i) / 10
		d2 := (i * i) % 10

		if (a[d1] && b[d2]) || (a[d2] && b[d1]) {
			//too lazy to negate this
		} else {
			return false
		}
	}
	return true
}

func trues() map[int]bool {
	die1 := make(map[int]bool)
	for i := 0; i < 10; i++ {
		die1[i] = true
	}

	return die1
}

func main() {
	starttime := time.Now()

	die1 := trues()

	die2 := trues()

	total := 0

	for i := 0; i < 10; i++ {
		for j := i + 1; j < 10; j++ {
			for k := j + 1; k < 10; k++ {
				for l := k + 1; l < 10; l++ {

					die1 = trues()
					die1[i] = false
					die1[j] = false
					die1[k] = false
					die1[l] = false
					if die1[6] || die1[9] {
						die1[9] = true
						die1[6] = true
					}

					for I := 0; I < 10; I++ {
						for J := I + 1; J < 10; J++ {
							for K := J + 1; K < 10; K++ {
								for L := K + 1; L < 10; L++ {

									die2 = trues()
									die2[I] = false
									die2[J] = false
									die2[K] = false
									die2[L] = false
									if die2[6] || die2[9] {
										die2[9] = true
										die2[6] = true
									}

									if works(die1, die2) {
										total++

									}

								}
							}
						}
					}
				}

			}
		}
	}
	//Divide by two for dice order
	fmt.Println(total / 2)

	fmt.Println("Elapsed time:", time.Since(starttime))
}
