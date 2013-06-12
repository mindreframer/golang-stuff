package main

import (
	"./euler"
	"fmt"
	"time"
)

func isPanindrome(n int64) bool {
	return n == euler.IntReverse(n)
}

const lid = 10000
const height = 100000000

func main() {
	starttime := time.Now()

	total := int64(0)
	dupes := map[int64]bool{}

	for i := int64(1); i < lid; i++ {
		sum := i * i
		sum += (i + 1) * (i + 1)

		for length := int64(2); sum < height; length++ {

			if isPanindrome(sum) && !dupes[sum] {
				dupes[sum] = true
				total += sum
			}

			sum += (i + length) * (i + length)
		}

	}

	fmt.Println(total)

	fmt.Println("Elapsed time:", time.Since(starttime))

}
