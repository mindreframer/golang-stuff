package main

import (
	"./euler"
	"fmt"
	"strconv"
	"time"
)

const zero = .001

//returns new list 
func apply(list []float64, place int, operation int) []float64 {
	answer := make([]float64, len(list)-1)
	copy(answer[:place], list[:place])
	answer[place] = operate(list[place], list[place+1], operation)
	copy(answer[place+1:], list[place+2:])
	return answer
}

func operate(x, y float64, which int) float64 {
	if which == 0 {
		return x + y
	} else if which == 1 {
		return x - y
	} else if which == 2 {
		return x * y
	} else if which == 3 {
		return x / y
	}
	return 0
}

func isInt(a float64) int {
	if a-float64(int(a)) > zero {
		return -1
	}
	return int(a)

}

func main() {
	starttime := time.Now()

	longestStreak := 0
	answer := ""

	for a := float64(1); a < 10; a++ {

		for b := a + 1; b < 10; b++ {
			for c := b + 1; c < 10; c++ {
				for d := c + 1; d < 10; d++ {

					ordered := []float64{a, b, c, d}

					table := make(map[int]bool)

					for reorder := 0; reorder < 24; reorder++ {

						for op1 := 0; op1 < 4; op1++ {
							for op2 := 0; op2 < 4; op2++ {
								for op3 := 0; op3 < 4; op3++ {
									for pl1 := 0; pl1 < 3; pl1++ {
										for pl2 := 0; pl2 < 2; pl2++ {
											list := make([]float64, len(ordered))
											copy(list, ordered)
											list = euler.PermuteFloats(reorder, list)
											list = apply(list, pl1, op1)
											list = apply(list, pl2, op2)
											list = apply(list, 0, op3)
											if isInt(list[0]) > 0 {
												table[isInt(list[0])] = true
											}

										}
									}
								}
							}
						}

					}

					best := 0
					streak := false
					streakcount := 0
					for i := 1; i <= 9*9*9*9; i++ {

						if table[i] {
							if streak {
								streakcount++
							}

							streak = true
						} else {
							streak = false
							if streakcount > best {
								best = streakcount
							}
							streakcount = 0
						}
					}

					if best > longestStreak {
						longestStreak = best
						//fmt.Println(a, b, c, d, ":", longestStreak)
						answer = strconv.Itoa(int(a)) + strconv.Itoa(int(b)) + strconv.Itoa(int(c)) + strconv.Itoa(int(d))
					}
				}
			}
		}
	}

	fmt.Println(answer)
	fmt.Println("Elapsed time:", time.Since(starttime))

}
