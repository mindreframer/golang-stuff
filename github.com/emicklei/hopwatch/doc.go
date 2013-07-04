// Copyright 2012 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
Hopwatch is a debugging tool for Go programs.  

Hopwatch uses a (embedded) HTML5 application to connect to your program (using a Websocket).
Using Hopwatch requires adding function calls at points of interest that allow you to watch program state and suspend the program.
On the Hopwatch page, you can view debug information (file:line,stack) and choose to resume the execution of your program.

You can provide more debug information using the Display and Dump functions which take an arbitrary number of variables.
The Display and Dump functions do not suspend the program ; it is like having logging information in the browser.

Usage:

	import (
		"github.com/emicklei/hopwatch"
	)

	func foo() {
		bar := "john"
		// suspends execution until hitting "Resume" in the browser
		hopwatch.Display("foo", bar).Break()
	}

Connect:

	The Hopwatch debugger is automatically started on http://localhost:23456/hopwatch.html.
	Your browser must support WebSockets. It has been tested with Chrome and Safari on a Mac.

Other code examples:

	// zero or more conditions ; conditionally suspends program (or goroutine)
	hopwatch.Break(i > 10,  j < 100)	

	// zero or more name,value pairs ; no program suspend
	hopwatch.Display("i",i , "j",j")

	// print any formatted string ; no program suspend
	hopwatch.Printf("result=%v", result)

	// display detailed (type, nesting) information using https://github.com/davecgh/go-spew
	hopwatch.Dump(myVar1)

	// format and display detailed (type, nesting) information using https://github.com/davecgh/go-spew
	hopwatch.Dumpf("myVar1: %v -- myVar2: %+v", myVar1, myVar2)	

Flags:

	-hopwatch	if set to false then hopwatch is disabled.
	-hopwatch.open	if set to false then hopwatch will not try to open the debugger page on startup.
	-hopwatch.break	if set to false then hopwatch will not suspend the program when Break(..) is called.
	-hopwatch.host	tcp hostname of the listener address (default = localhost).
	-hopwatch.port	tcp port of the listener address (default = 23456).

Install:

	go get -u github.com/emicklei/hopwatch


Resources:

	https://github.com/emicklei/hopwatch (project)
	http://ernestmicklei.com/2012/12/14/hopwatch-a-debugging-tool-for-go/  (blog)


(c) 2012-2013, Ernest Micklei. MIT License
*/
package hopwatch
