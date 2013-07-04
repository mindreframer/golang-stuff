package main

import (
	"github.com/emicklei/hopwatch"
	"time"
	"os"
)

func main() {
	haveABreak()
}

func haveABreak() {
	hopwatch.Break()
	printNow()	
}

func printNow() {
	hopwatch.Printf("time is: %v", time.Now())
	dumpArgs()
}

func dumpArgs() {
	hopwatch.Dump(os.Args).Break()
	waitHere()
}