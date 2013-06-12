package main

import (
	"euler"
	"fmt"
	"time"
)

func main() {
	starttime := time.Now()

	fmt.Println("Hello, World", euler.Prime(10000))

	fmt.Println("Elapsed time:", time.Since(starttime))
}
