package main

import (
	"github.com/emicklei/hopwatch"
	"math"
)

func main() {
	liveOfPi()
}

func liveOfPi() {
	hopwatch.Dump(math.Pi).Break()
}
