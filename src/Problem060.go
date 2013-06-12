package main

import (
	"./euler"
	"fmt"
)

func works(i int64, j int64) bool {

	if !euler.IsPrime(euler.ConcatanInt(euler.Prime(i), euler.Prime(j))) {
		return false
	}
	if !euler.IsPrime(euler.ConcatanInt(euler.Prime(j), euler.Prime(i))) {
		return false
	}

	return true
}

func main() {

	ceiling := int64(2000)
	current := int64(99999999999)

	for i := int64(1); i < ceiling; i++ {

		for j := i; j < ceiling; j++ {

			if works(i, j) {

				for k := j; k < ceiling; k++ {

					if works(i, k) && works(j, k) {

						for l := k; l < ceiling; l++ {

							if works(i, l) && works(j, l) && works(k, l) {

								for a := l; a < ceiling; a++ {
									if works(a, i) && works(a, j) && works(a, k) && works(a, l) {
										sum := euler.Prime(a) + euler.Prime(i) + euler.Prime(j) + euler.Prime(k) + euler.Prime(l)
										if sum < current {
											fmt.Println(euler.Prime(i), euler.Prime(j), euler.Prime(k), euler.Prime(l), euler.Prime(a), "Sum:", sum)
											current = sum

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
