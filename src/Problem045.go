package main

import "fmt"

var current3 int64 = 0
var current5 int64 = 0
var current6 int64 = 0

func triangle() int64 {

	current3++
	return current3 * (current3 + 1) / 2
}

func pentagon() int64 {

	current5++
	return current5 * (3*current5 - 1) / 2
}

func hexagon() int64 {

	current6++
	return current6 * (2*current6 - 1)
}

func main() {

//	var counter int64 = 2

	var pent int64 =0
	var hex int64 =0
	var tri int64 =0

	for {

		for pent < tri || pent < hex{
			pent = pentagon()
		}
		for tri < pent || tri < hex{
			tri = triangle()
		}
		for hex < pent || hex < tri{
			hex = hexagon()
		}

		if tri == pent && tri == hex {
			fmt.Println(tri,":",current3, current5, current6)
			tri = triangle()
		}
		
//		if tri > counter {
//			fmt.Println("We're at", tri, "(triangle number ",current3,")")
//			counter = int64 (2* float64(counter))
//		}	
		
	}
}
