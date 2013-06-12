package main

import (
	"fmt"
)

func istripple(a int, b int) bool {
	max := b
	if a > b {
		max = a
	}

	target := int64(a)*int64(a) + int64(b)*int64(b)

	c := int64(max)

	for c*c < target {

		c++
	}
	return c*c == target
}

func third(a int, b int) int64 {
	max := b
	if a > b {
		max = a
	}

	target := int64(a)*int64(a) + int64(b)*int64(b)

	c := int64(max)

	for c*c < target {

		c++
	}
	return c
}

func main() {

	total := int64(0)
	count := 0

	b := 2

	for count < 12 {
		for offset := -1; offset < 3; offset += 2 {

			h := b + offset
			if istripple(b, 2*h) && third(b, 2*h)%2 == 0 {
				l := third(b, 2*h) / 2
				fmt.Println(b, h, l)
				total += l
				count++
			}

		}
		b++

	}

	fmt.Println(total)
}
