package main

import (
	"euler"
	"fmt"
	"math/big"
	"time"
)

const (
	lid    = 12
	target = 4000000
	top    = 4
)

func evaluate(factors [][2]int64) *big.Int {
	show := big.NewInt(1)
	for i := 0; i < len(factors); i++ {
		if factors[i][1] > 0 {
			prime := big.NewInt(factors[i][0])

			for j := int64(0); j < factors[i][1]; j++ {
				show = show.Mul(show, prime)
			}
		}
	}
	return show
}

func solns(factors [][2]int64) int64 {
	answer := int64(1)

	for i := 0; i < len(factors); i++ {
		answer *= 1 + (2 * factors[i][1])
	}

	return answer/2 + 1
}

func main() {
	starttime := time.Now()

	fac := make([][2]int64, lid)
	for i := int64(1); i < lid+1; i++ {
		fac[i-1] = [2]int64{euler.Prime(i), 0}
	}

	var smallest big.Int
	smallest.SetString("9999999999999999", 10)

	for fac[0][1] = 0; fac[0][1] < top; fac[0][1]++ {
		for fac[1][1] = fac[0][1]; fac[1][1] >= 0; fac[1][1]-- {
			for fac[2][1] = fac[1][1]; fac[2][1] >= 0; fac[2][1]-- {
				for fac[3][1] = fac[2][1]; fac[3][1] >= 0; fac[3][1]-- {
					for fac[4][1] = fac[3][1]; fac[4][1] >= 0; fac[4][1]-- {
						for fac[5][1] = fac[4][1]; fac[5][1] >= 0; fac[5][1]-- {
							for fac[6][1] = fac[5][1]; fac[6][1] >= 0; fac[6][1]-- {
								for fac[7][1] = fac[6][1]; fac[7][1] >= 0; fac[7][1]-- {
									for fac[8][1] = fac[7][1]; fac[8][1] >= 0; fac[8][1]-- {
										for fac[9][1] = fac[8][1]; fac[9][1] >= 0; fac[9][1]-- {
											for fac[10][1] = fac[9][1]; fac[10][1] >= 0; fac[10][1]-- {
												for fac[11][1] = fac[10][1]; fac[11][1] >= 0; fac[11][1]-- {
													if solns(fac) > target && smallest.Cmp(evaluate(fac)) > 0 {
														smallest = *evaluate(fac)
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	fmt.Println(smallest.String())

	fmt.Println("Elapsed time:", time.Since(starttime))
}
