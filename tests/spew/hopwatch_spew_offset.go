package main

import (
	"github.com/emicklei/hopwatch"
)

func main() {
	hopwatch.Dump(8).Dump(8).Break()
	hopwatch.Dumpf("%v", 9).Dumpf("%v", 9).Break()
}
