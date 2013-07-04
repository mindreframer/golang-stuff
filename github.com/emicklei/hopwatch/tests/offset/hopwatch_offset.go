package main

import (
	"github.com/emicklei/hopwatch"
)

func main() {
	hopwatch.Display("8",8)
	hopwatch.Display("9",9).Break()
	inside()
	indirectDisplay("11",11)
	indirectBreak()
	illegalOffset()
}
func inside() {
	hopwatch.Display("16",16)
	hopwatch.Display("17",17).Break()
}
func indirectDisplay(args ...interface{}) {
	hopwatch.CallerOffset(2).Display(args...)
}
func indirectBreak() {
	hopwatch.CallerOffset(3).Break()
}
func illegalOffset() {
	defer func() {
		if r := recover(); r != nil {
            print("Recovered in illegalOffset")
        }
	}()
	hopwatch.CallerOffset(-1).Break()
}
