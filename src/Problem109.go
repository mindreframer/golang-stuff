package main

import (
	"fmt"
	"time"
)

const target = 100

func main() {
	starttime := time.Now()
	counter := 0

	for i := 0; i <= 21; i++ {

		for j := i; j <= 21; j++ {
			if i == 21 {
				i = 25
			}
			if j == 21 {
				j = 25
			}

			for k := 1; k <= 21; k++ {
				if k == 21 {
					k = 25
				}

				for a := 1; a <= 3 && !(i == 0 && a > 1); a++ {
					for b := a; b <= 3 && !(j == 0 && b > 1); b++ {
						if !((b == 3 && j == 25) || (a == 3 && i == 25)) {

							if (a*i)+(b*j)+(2*k) < target {

								counter++
							}
						}
					}
				}

			}
		}
	}

	for i := 1; i <= 21; i++ {
		if i == 21 {
			i = 25
		}

		for j := i + 1; j <= 21; j++ {

			if j == 21 {
				j = 25
			}

			for k := 1; k <= 21; k++ {
				if k == 21 {
					k = 25
				}
				for b := 1; b <= 3; b++ {

					for a := b + 1; a <= 3; a++ {
						if !((b == 3 && j == 25) || (a == 3 && i == 25)) {
							if (a*i)+(b*j)+(2*k) < target {

								counter++
							}

						}

					}
				}

			}
		}
	}

	fmt.Println(counter)

	fmt.Println("Elapsed time:", time.Since(starttime))
}
