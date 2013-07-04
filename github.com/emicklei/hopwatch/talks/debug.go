package main

import (
	"github.com/emicklei/hopwatch"
)

func debug(vars ...interface{}) *hopwatch.Watchpoint {

	// create a watchpoint and compensate for extra stack frame
	watchPoint := hopwatch.CallerOffset(2 + 1)

	watchPoint.Dump(vars...)

	// revert compensation for calls on the return value
	watchPoint.CallerOffset(2)

	return watchPoint
}
