// Copyright 2012,2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package hopwatch

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
)

// command is used to transport message to and from the debugger.
type command struct {
	Action     string
	Parameters map[string]string
}

// addParam adds a key,value string pair to the command ; no check on overwrites.
func (self *command) addParam(key, value string) {
	if self.Parameters == nil {
		self.Parameters = map[string]string{}
	}
	self.Parameters[key] = value
}

var (
	hopwatchHostParam  = flag.String("hopwatch.host", "localhost", "HTTP host the debugger is listening on")
	hopwatchPortParam  = flag.Int("hopwatch.port", 23456, "HTTP port the debugger is listening on")
	hopwatchParam      = flag.Bool("hopwatch", true, "controls whether hopwatch agent is started")
	hopwatchOpenParam  = flag.Bool("hopwatch.open", true, "controls whether a browser page is opened on the hopwatch page")
	hopwatchBreakParam = flag.Bool("hopwatch.break", true, "do not suspend the program if Break(..) is called")

	hopwatchEnabled            = true
	hopwatchOpenEnabled        = true
	hopwatchBreakEnabled       = true
	hopwatchHost               = "localhost"
	hopwatchPort         int64 = 23456

	currentWebsocket   *websocket.Conn
	toBrowserChannel   = make(chan command)
	fromBrowserChannel = make(chan command)
	connectChannel     = make(chan command)
	debuggerMutex      = sync.Mutex{}
)

func init() {
	// check any command line params. (needed when programs do not call flag.Parse() )
	for i, arg := range os.Args {
		if strings.HasPrefix(arg, "-hopwatch=") {
			if strings.HasSuffix(arg, "false") {
				log.Printf("[hopwatch] disabled.\n")
				hopwatchEnabled = false
				return
			}
		}
		if strings.HasPrefix(arg, "-hopwatch.open") {
			if strings.HasSuffix(arg, "false") {
				log.Printf("[hopwatch] auto open debugger disabled.\n")
				hopwatchOpenEnabled = false
			}
		}
		if strings.HasPrefix(arg, "-hopwatch.break") {
			if strings.HasSuffix(arg, "false") {
				log.Printf("[hopwatch] suspend on Break(..) disabled.\n")
				hopwatchBreakEnabled = false
			}
		}
		if strings.HasPrefix(arg, "-hopwatch.host") {
			if eq := strings.Index(arg, "="); eq != -1 {
				hopwatchHost = arg[eq+1:]
			} else if i < len(os.Args) {
				hopwatchHost = os.Args[i+1]
			}
		}
		if strings.HasPrefix(arg, "-hopwatch.port") {
			portString := ""
			if eq := strings.Index(arg, "="); eq != -1 {
				portString = arg[eq+1:]
			} else if i < len(os.Args) {
				portString = os.Args[i+1]
			}
			port, err := strconv.ParseInt(portString, 10, 8)
			if err != nil {
				log.Panicf("[hopwatch] illegal port parameter:%v", err)
			}
			hopwatchPort = port
		}
	}
	http.HandleFunc("/hopwatch.html", html)
	http.HandleFunc("/hopwatch.css", css)
	http.HandleFunc("/hopwatch.js", js)
	http.HandleFunc("/gosource", gosource)
	http.Handle("/hopwatch", websocket.Handler(connectHandler))
	go listen()
	go sendLoop()
}

// Open calls the OS default program for uri
func open(uri string) error {
	var run string
	switch {
	case "windows" == runtime.GOOS:
		run = "start"
	case "darwin" == runtime.GOOS:
		run = "open"
	case "linux" == runtime.GOOS:
		run = "xdg-open"
	default:
		return fmt.Errorf("Unable to open uri:%v on:%v", uri, runtime.GOOS)
	}
	return exec.Command(run, uri).Start()
}

// serve a (source) file for displaying in the debugger
func gosource(w http.ResponseWriter, req *http.Request) {
	fileName := req.FormValue("file")
	// should check for permission?  
	w.Header().Set("Cache-control", "no-store, no-cache, must-revalidate")
	http.ServeFile(w, req, fileName)
}

// listen starts a Http Server on a fixed port.
// listen is run in parallel to the initialization process such that it does not block.
func listen() {
	hostPort := fmt.Sprintf("%s:%d", hopwatchHost, hopwatchPort)
	if hopwatchOpenEnabled {
		log.Printf("[hopwatch] opening http://%v/hopwatch.html ...\n", hostPort)
		go open(fmt.Sprintf("http://%v/hopwatch.html", hostPort))
	} else {
		log.Printf("[hopwatch] open http://%v/hopwatch.html ...\n", hostPort)
	}
	if err := http.ListenAndServe(hostPort, nil); err != nil {
		log.Printf("[hopwatch] failed to start listener:%v", err.Error())
	}
}

// connectHandler is a Http handler and is called on loading the debugger in a browser.
// As soon as a command is received the receiveLoop is started. 
func connectHandler(ws *websocket.Conn) {
	if currentWebsocket != nil {
		log.Printf("[hopwatch] already connected to a debugger; Ignore this\n")
		return
	}
	// remember the connection for the sendLoop	
	currentWebsocket = ws
	var cmd command
	if err := websocket.JSON.Receive(currentWebsocket, &cmd); err != nil {
		log.Printf("[hopwatch] connectHandler.JSON.Receive failed:%v", err)
	} else {
		log.Printf("[hopwatch] connected to browser. ready to hop")
		connectChannel <- cmd
		receiveLoop()
	}
}

// receiveLoop reads commands from the websocket and puts them onto a channel.
func receiveLoop() {
	for {
		var cmd command
		if err := websocket.JSON.Receive(currentWebsocket, &cmd); err != nil {
			log.Printf("[hopwatch] receiveLoop.JSON.Receive failed:%v", err)
			fromBrowserChannel <- command{Action: "quit"}
			break
		}
		if "quit" == cmd.Action {
			hopwatchEnabled = false
			log.Printf("[hopwatch] browser requests disconnect.\n")
			currentWebsocket.Close()
			currentWebsocket = nil
			fromBrowserChannel <- cmd
			break
		} else {
			fromBrowserChannel <- cmd
		}
	}
}

// sendLoop takes commands from a channel to send to the browser (debugger).
// If no connection is available then wait for it.
// If the command action is quit then abort the loop.
func sendLoop() {
	if currentWebsocket == nil {
		log.Print("[hopwatch] no browser connection, wait for it ...")
		cmd := <-connectChannel
		if "quit" == cmd.Action {
			return
		}
	}
	for {
		next := <-toBrowserChannel
		if "quit" == next.Action {
			break
		}
		if currentWebsocket == nil {
			log.Print("[hopwatch] no browser connection, wait for it ...")
			cmd := <-connectChannel
			if "quit" == cmd.Action {
				break
			}
		}
		websocket.JSON.Send(currentWebsocket, &next)
	}
}

// watchpoint is a helper to provide a fluent style api.
// This allows for statements like hopwatch.Display("var",value).Break()
type Watchpoint struct {
	disabled bool
	offset   int // offset in the caller stack for highlighting source
}

// Printf formats according to a format specifier and writes to the debugger screen. 
// It returns a new Watchpoint to send more or break.
func Printf(format string, params ...interface{}) *Watchpoint {
	wp := &Watchpoint{offset: 2}
	return wp.Printf(format, params...)
}

// Display sends variable name,value pairs to the debugger.
// The parameter nameValuePairs must be even sized.
func Display(nameValuePairs ...interface{}) *Watchpoint {
	wp := &Watchpoint{offset: 2}
	return wp.Display(nameValuePairs...)
}

// Break suspends the execution of the program and waits for an instruction from the debugger (e.g. Resume).
// Break is only effective if all (if any) conditions are true. The program will resume otherwise.
func Break(conditions ...bool) {
	suspend(2, conditions...)
}

// CallerOffset (default=2) allows you to change the file indicator in hopwatch.
// Use this method when you wrap the .CallerOffset(..).Display(..).Break() in your own function.
func CallerOffset(offset int) *Watchpoint {
	return (&Watchpoint{}).CallerOffset(offset)
}

// CallerOffset (default=2) allows you to change the file indicator in hopwatch.
func (w *Watchpoint) CallerOffset(offset int) *Watchpoint {
	if hopwatchEnabled && (offset < 0) {
		log.Panicf("[hopwatch] ERROR: illegal caller offset:%v . watchpoint is disabled.\n", offset)
		w.disabled = true
	}
	w.offset = offset
	return w
}

// Printf formats according to a format specifier and writes to the debugger screen. 
func (self *Watchpoint) Printf(format string, params ...interface{}) *Watchpoint {
	self.offset += 1
	var content string
	if len(params) == 0 {
		content = format
	} else {
		content = fmt.Sprintf(format, params...)
	}
	return self.printcontent(content)
}

// Printf formats according to a format specifier and writes to the debugger screen. 
func (self *Watchpoint) printcontent(content string) *Watchpoint {
	_, file, line, ok := runtime.Caller(self.offset)
	cmd := command{Action: "print"}
	if ok {
		cmd.addParam("go.file", file)
		cmd.addParam("go.line", fmt.Sprint(line))
	}
	cmd.addParam("line", content)
	channelExchangeCommands(cmd)
	return self
}

// Display sends variable name,value pairs to the debugger. Values are formatted using %#v.
// The parameter nameValuePairs must be even sized.
func (self *Watchpoint) Display(nameValuePairs ...interface{}) *Watchpoint {
	_, file, line, ok := runtime.Caller(self.offset)
	cmd := command{Action: "display"}
	if ok {
		cmd.addParam("go.file", file)
		cmd.addParam("go.line", fmt.Sprint(line))
	}
	if len(nameValuePairs)%2 == 0 {
		for i := 0; i < len(nameValuePairs); i += 2 {
			k := nameValuePairs[i]
			v := nameValuePairs[i+1]
			cmd.addParam(fmt.Sprint(k), fmt.Sprintf("%#v", v))
		}
	} else {
		log.Printf("[hopwatch] WARN: missing variable for Display(...) in: %v:%v\n", file, line)
		self.disabled = true
		return self
	}
	channelExchangeCommands(cmd)
	return self
}

// Break halts the execution of the program and waits for an instruction from the debugger (e.g. Resume).
// Break is only effective if all (if any) conditions are true. The program will resume otherwise.
func (self Watchpoint) Break(conditions ...bool) {
	suspend(self.offset, conditions...)
}

// suspend will create a new Command and send it to the browser.
// callerOffset controls from which stackframe the go source file and linenumber must be read.
// Ignore if option hopwatch.break=false
func suspend(callerOffset int, conditions ...bool) {
	if !hopwatchBreakEnabled {
		return
	}
	for _, condition := range conditions {
		if !condition {
			return
		}
	}
	_, file, line, ok := runtime.Caller(callerOffset)
	cmd := command{Action: "break"}
	if ok {
		cmd.addParam("go.file", file)
		cmd.addParam("go.line", fmt.Sprint(line))
		cmd.addParam("go.stack", trimStack(string(debug.Stack())))
	}
	channelExchangeCommands(cmd)
}

// Peel off the part of the stack that lives in hopwatch
func trimStack(stack string) string {
	lines := strings.Split(stack, "\n")
	c := 0
	for _, line := range lines {
		if strings.Index(line, "/hopwatch") == -1 { // means no function in this package
			break
		}
		c++
	}
	return strings.Join(lines[c:], "\n")
}

// Put a command on the browser channel and wait for the reply command
func channelExchangeCommands(toCmd command) {
	if !hopwatchEnabled {
		return
	}
	// synchronize command exchange ; break only one goroutine at a time
	debuggerMutex.Lock()
	toBrowserChannel <- toCmd
	<-fromBrowserChannel
	debuggerMutex.Unlock()
}
