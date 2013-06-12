package main

import (
	"fmt"
	"time"
)

var memo = make(map[int64][2]int64)

const (
	target = 1000000000000
	mod    = 100000
)

func fac(n int64) [2]int64 {
	if answer, ok := memo[n]; ok {
		return answer
	}

	base := fac(n - 1)[1]

	answer := fac(n - 1)[0]

	offset := int64(0)
	for n%10 == 0 {
		n = n / 10
		offset++

	}

	answer *= n

	for answer%10 == 0 {
		answer = answer / 10
		offset++
	}

	answer = answer % mod

	memo[n] = [2]int64{answer, base + offset}

	//hogs memory otherwise
	delete(memo, n-5)

	return [2]int64{answer, base + offset}

}

func prod(a, b [2]int64) [2]int64 {

	offset := a[1] + b[1]
	for a[0]%10 == 0 {
		a[0] = a[0] / 10
		offset++
	}
	for b[0]%10 == 0 {
		b[0] = b[0] / 10
		offset++
	}

	return [2]int64{(a[0] * b[0]) % mod, offset}
}

func exp(a [2]int64, pow int64) [2]int64 {
	if pow == 1 {
		return a
	}
	if pow%2 == 0 {
		return prod(exp(a, pow/2), exp(a, pow/2))
	}
	return prod(prod(exp(a, pow/2), exp(a, pow/2)), a)
}

func main() {
	starttime := time.Now()

	memo[0] = [2]int64{1, 0}

	fmt.Println(fac(mod - 1))

	intermediate := fac(mod - 1)

	fmt.Println(exp(intermediate, target/mod))

	fmt.Println("Elapsed time:", time.Since(starttime))
}
