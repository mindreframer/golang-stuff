package main

import "fmt"

func main() {
	record := 0

	for perimeter := 0; perimeter < 1000; perimeter++ {

		count := 0

		for a := 1; a < perimeter; a++ {

			for b := a + 1; b < perimeter; b++ {

				for c := b + 1; c < perimeter; c++ {

					if a+b+c == perimeter && a*a+b*b == c*c {
						count++

					}

				}

			}

		}

		if count > record {
			record = count
			fmt.Println(record, perimeter)

		}

	}

}
