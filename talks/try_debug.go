package main

import "math"

func main() {
	liveOfPi()
}

func liveOfPi() {
	debug(math.Pi).Break()
}
