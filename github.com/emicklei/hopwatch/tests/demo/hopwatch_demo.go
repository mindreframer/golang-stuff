package main

import (
	"github.com/emicklei/hopwatch"
)

func main() {
	for i := 0; i < 6; i++ {
		hopwatch.Display("i", i)
		j := i * i
		hopwatch.Display("i", i, "j", j).Break(j > 10)
		hopwatch.Printf("%#v", "printf formatted value(s)")
		hopwatch.Break()
		quick()
	}
}

func quick() {
	hopwatch.Break()
	brown()
}
func brown() {
	hopwatch.Break()
	fox()
}
func fox() {
	hopwatch.Break()
	jumps()
}
func jumps() {
	hopwatch.Break()
	over()
}
func over() {
	hopwatch.Break()
	the()
}
func the() {
	hopwatch.Break()
	lazy()
}
func lazy() {
	hopwatch.Break()
	dog()
}
func dog() {
	hopwatch.Break()
}
