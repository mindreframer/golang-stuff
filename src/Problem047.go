package main

import (
	"./eulerlib"
	"fmt"
)

const length = 4

func works(start int64) bool {
	for i := int64(0); i < length; i++ {
		if eulerlib.DistinctNumber(eulerlib.Factor(start+i)) != length {
			return false
		}
	}
	return true
}

func main() {

	for i := int64(0); i < 150000; i++ {
		if works(i) {
			fmt.Println(i)
		}
	}

}
