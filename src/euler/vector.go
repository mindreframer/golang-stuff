package euler

import (
//"fmt"
)

func BubbleSortVec(input []string) {

	for j := len(input) - 1; j > 0; j-- {

		for i := 0; i < j; i++ {
			if input[i] > input[i+1] {
				temp := input[i]
				input[i] = input[i+1]
				input[i+1] = temp

			}
		}

	}

}

//remove entry i to j inclusive (modified from 
//(http://golang.org/doc/articles/slices_usage_and_internals.html?h=slice+pointer)3
func CutVec(a []string, i int, j int) {
	copy(a[i:], a[j+1:])
	for k, n := len(a)-j-1+i, len(a); k < n; k++ {
		a[k] = "" // or the zero value of T
	} // for k
	a = a[:len(a)-j-1+i]
}

//removes duplicate entries in a SORTED vector.
func RemoveDuplicatesVec(input []string) []string {

	if len(input) == 1 {
		return input
	}

	if input[0] == input[1] {
		return RemoveDuplicatesVec(append(input[1:]))
	}

	return append(input[:1], RemoveDuplicatesVec(input[1:])...)

}
