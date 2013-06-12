package main

import (
	"./euler"
	"fmt"
	"time"
)

const (
	tablesize = 10000000
	tailength = 1000000000
	zero      = .00001
)

var (
	fibTails [tablesize]int
	fibHeads [tablesize]science
)

//a la "scientific notation"
type science struct {
	value  float64
	offset int
}

func fibTail(n int) int {

	if n <= 2 {
		return 1
	}

	if n < tablesize && fibTails[n] != 0 {
		return fibTails[n]
	}

	answer := (fibTail(n-1) + fibTail(n-2)) % tailength
	if n < tablesize {
		fibTails[n] = answer
	}
	return answer

}

func fibHead(n int) science {

	if n <= 2 {
		return science{1, 0}
	}

	if n < tablesize && fibHeads[n].value > zero {
		return fibHeads[n]
	}

	part1 := fibHead(n - 1)
	part2 := fibHead(n - 2)

	var answer science

	if part1.offset == part2.offset {
		answer.value = part1.value + part2.value
		answer.offset = part1.offset
	}

	if part1.offset > part2.offset {
		answer.value = part1.value + (part2.value / 10)
		answer.offset = part1.offset
	}

	if answer.value > tailength {
		answer.value = answer.value / 10
		answer.offset++
	}

	if n < tablesize {
		fibHeads[n] = answer
	}
	return answer

}

func main() {
	starttime := time.Now()

	answer := 0

	for i := 10; !euler.IsPandigital(int64(fibHead(i).value)) || !euler.IsPandigital(int64(fibTail(i))); i++ {
		answer = i + 1
	}

	fmt.Println(answer)

	fmt.Println("Elapsed time:", time.Since(starttime))

}
