package main

import (
	//"./euler"
	"./euler"
	"fmt"
	"math"
	"time"
)

//fraction of the form (a +  b \sqrt R) / d
type frac struct {
	a int
	b int
	R int
	d int
}

func flip(F frac) frac {
	return frac{
		a: F.a * F.d,
		b: -F.b * F.d,
		R: F.R,
		d: (F.a * F.a) - (F.b * F.b * F.R)}
}

func reduce(F frac) frac {
	gcd := int(euler.GCD(int64(F.a), int64(F.b)))
	gcd = int(euler.GCD(int64(gcd), int64(F.d)))
	return frac{F.a / gcd, F.b / gcd, F.R, F.d / gcd}
}

func nextFrac(F frac) (n int, next frac) {
	total := (float64(F.a) + (float64(F.b) * math.Sqrt(float64(F.R)))) / float64(F.d)

	n = int(total)

	if n != 0 {
		next = frac{
			a: F.a - (n * F.d),
			b: F.b,
			R: F.R,
			d: F.d}
	} else {
		n = 1
		next = frac{
			a: F.a + F.d,
			b: F.b,
			R: F.R,
			d: F.d}
	}

	next = flip(next)
	next = reduce(next)

	return
}

func isSquare(n int) bool {
	sqrt := int(math.Sqrt(float64(n)))
	if sqrt*sqrt == n {
		return true
	}
	return false
}

func main() {
	starttime := time.Now()

	oddPeriod := 0

	for rad := 2; rad <= 10000; rad++ {

		if !isSquare(rad) {

			test := frac{a: 0, b: 1, R: rad, d: 1}
			n := 0

			//Clear any aperiodic part of the expansion
			for i := 0; i < 10; i++ {
				n, test = nextFrac(test)

			}

			start, startfrac := n, test
			length := 1
			n, test = nextFrac(test)
			for n != start || test != startfrac {
				n, test = nextFrac(test)
				length++
			}

			if length%2 == 1 {
				oddPeriod++
			}

			//fmt.Println(rad, "has period", length)
		}

	}

	fmt.Println(oddPeriod)

	fmt.Println("Elapsed time:", time.Since(starttime))

}
