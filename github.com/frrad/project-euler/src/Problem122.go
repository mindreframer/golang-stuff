package main

import (
	"fmt"
	"time"
)

type tree struct {
	parent *tree
	value  int
}

const searchLength = 200 + 1

func main() {
	starttime := time.Now()

	table := [searchLength]int{}

	haveseen := make(map[int]bool)
	haveseen[1] = true
	table[1] = 0

	current := make([]*tree, 1)
	current[0] = &tree{value: 1}

	for i := 0; i < 14; i++ {

		next := make([]*tree, 0)
		for _, leaf := range current {
			n := leaf.value

			ancestor := leaf
			temp := make([]*tree, 0)
			for ancestor != nil {
				consider := n + ancestor.value
				if !haveseen[consider] {

					if consider < searchLength {
						table[consider] = i + 1

					}
					haveseen[consider] = true
					temp = append([]*tree{&tree{leaf, consider}}, temp...)
				}
				ancestor = ancestor.parent
			}
			next = append(next, temp...)
		}

		current = next

	}

	total := 0
	for i := 1; i < searchLength; i++ {
		total += table[i]
	}

	//Power method does not work for 2 values < 200 (see Knuth Vol 2. 4.6.3)
	fmt.Println(total - 2)

	fmt.Println("Elapsed time:", time.Since(starttime))

}
