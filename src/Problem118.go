package main

import (
	"./euler"
	"fmt"
	"strconv"
	"time"
)

func paste(a, b int) (answer int) {
	answer, _ = strconv.Atoi(strconv.Itoa(a) + strconv.Itoa(b))
	return
}

func glue(n int, list []int) []int {

	if n == 0 || len(list) == 1 {
		return list
	}

	if n%2 == 0 {
		return append(glue(n/2, list[:len(list)-1]), list[len(list)-1])
	}

	mist := paste(list[len(list)-2], list[len(list)-1])
	return glue(n/2, append(list[:len(list)-2], mist))

}

func main() {
	starttime := time.Now()

	total := make([]int, 10)

	for j := 0; j < 362880; j++ {

		for i := 0; i < 256; i++ {

			current := euler.Permutation(j, []int{1, 2, 3, 4, 5, 6, 7, 8, 9})
			current = glue(i, current)

			flag := true
			for k := 0; k < len(current) && flag == true; k++ {
				if !euler.IsPrime(int64(current[k])) {
					flag = false
				}
			}
			if flag == true {
				//fmt.Println(current)
				total[len(current)]++

			}

		}

		//		fmt.Println(j, total)
	}

	fmt.Println(total)

	fmt.Println(total[1] + (total[2] / 2) + (total[3] / 6) + (total[4] / 24) + (total[5] / 120) + (total[6] / 720) + (total[7] / 5040))

	fmt.Println("Elapsed time:", time.Since(starttime))

}
