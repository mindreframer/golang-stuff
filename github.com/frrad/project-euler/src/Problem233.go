package main

import (
	"./euler"
	"fmt"
	"strconv"
	"time"
)

func isSmall(str string) bool {
	if len(str) > 11 {
		return false
	}
	return true
}
func isMultiple(n int64) bool { //returns true if there are no pythagorean prime factors
	for _, factor := range euler.Factor(n) {
		if factor%4 == 1 {
			return false
		}
	}
	return true
}
func lastSmall(table []string, multiple string) int {
	min := 0
	max := len(table)
	for max-min > 1 {
		if isSmall(euler.StringProd(multiple, table[(min+max)/2])) {
			min = (min + max) / 2
		} else {
			max = (min + max) / 2
		}
	}
	return min
}
func main() {
	starttime := time.Now()

	primes1 := make([]string, 0) //Table of pythagorean primes
	for i := int64(1); i < 333000; i++ {
		num := euler.Prime(i)
		if num%4 == 1 {
			primes1 = append(primes1, strconv.FormatInt(num, 10))
		}
	}

	multitable := make([]string, 1) //numbers of the form 2^k \prod pi
	multitable[0] = "1"
	for i := int64(1); i < 280000; i++ {
		if isMultiple(i) {
			multitable = append(multitable, strconv.FormatInt(i, 10))
		}
	}

	sumtable := make([]string, len(multitable)) //partial sums of multitable
	sumtable[0] = multitable[0]
	for i := 1; i < len(multitable); i++ {
		sumtable[i] = euler.StringSum(sumtable[i-1], multitable[i])
	}

	total := "0"

	//primes of the form p1^p2*10^2 for pi distinct, pythagorean
	prime1 := ""
	for index1 := 0; isSmall(prime1); index1++ {
		prime1 = euler.StringExp(primes1[index1], 10)
		product := prime1
		for index2 := 0; isSmall(product); index2++ {
			if index2 == index1 {
				index2++
			}
			product = euler.StringProd(prime1, euler.StringExp(primes1[index2], 2))
			if isSmall(product) {
				contribution := euler.StringProd(product, sumtable[lastSmall(multitable, product)])
				total = euler.StringSum(total, contribution)
			}
		}
	}

	//primes of the form p1^7*p2^3 for pi distinct, pythagorean
	prime1 = ""
	for index1 := 0; isSmall(prime1); index1++ {
		prime1 = euler.StringExp(primes1[index1], 7)
		product := prime1
		for index2 := 0; isSmall(product); index2++ {
			if index2 == index1 {
				index2++
			}
			product = euler.StringProd(prime1, euler.StringExp(primes1[index2], 3))
			if isSmall(product) {
				contribution := euler.StringProd(product, sumtable[lastSmall(multitable, product)])
				total = euler.StringSum(total, contribution)
			}
		}
	}

	//primes of the form p1^3*p2^2*p3 for pi distinct, pythagorean
	prime1 = ""
	for index1 := 0; isSmall(prime1); index1++ { //Index of q3
		prime1 = euler.StringExp(primes1[index1], 3)
		product1 := prime1
		for index2 := 0; isSmall(product1); index2++ { //Index of q2
			if index2 == index1 {
				index2++
			}
			product1 = euler.StringProd(prime1, euler.StringExp(primes1[index2], 2))
			finalproduct := product1
			for index3 := 0; isSmall(finalproduct); index3++ { //Index of q3
				for index1 == index3 || index3 == index2 {
					index3++
				}
				finalproduct = euler.StringProd(product1, primes1[index3])
				if isSmall(finalproduct) {
					contribution := euler.StringProd(finalproduct, sumtable[lastSmall(multitable, finalproduct)])
					total = euler.StringSum(total, contribution)
				}
			}
		}
	}

	fmt.Println(total)
	fmt.Println("Elapsed time:", time.Since(starttime))
}
