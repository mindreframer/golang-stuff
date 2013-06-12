package main

import (
	"fmt"
)

func last10Of(height int) int64 {
	current := int64(1)
	for i := 0; i < height; i++ {
		current = (current * 2) % 10000000000
	}
	return current
}
func main() {
	fmt.Println((last10Of(7830457)*28433 + 1) % 10000000000)
}
