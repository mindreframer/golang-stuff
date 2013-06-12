// Copyright (C) 2012 Numerotron Inc.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package bingo

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var AccessLogFilename = "/tmp/access.log"
var ErrorLogFilename = "/tmp/error.log"

var alf *os.File
var elf *os.File

var accessLog *log.Logger
var errorLog *log.Logger

var qaccess chan string
var qerror chan string

func init() {
	qaccess = make(chan string, 1000)
	qerror = make(chan string, 1000)
	var err error
	alf, err = os.OpenFile(AccessLogFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		panic(fmt.Sprintf("couldn't open access log file: %s", err))
	}
	accessLog = log.New(alf, "", log.LstdFlags)
	elf, err = os.OpenFile(ErrorLogFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		panic("couldn't open error log file")
	}
	errorLog = log.New(elf, "", log.LstdFlags)

	go writeAccess()
	go writeErrors()
}

func LogAccess(req *http.Request, elapsed time.Duration) {
	// qaccess <- fmt.Sprintf("%s \"%s %s %s\" %dms", req.RemoteAddr, req.Method, req.RequestURI, req.Proto, elapsed / time.Millisecond)
	qaccess <- fmt.Sprintf("%s \"%s %s %s\" %s", req.RemoteAddr, req.Method, req.RequestURI, req.Proto, elapsed)
}

func LogError(req *http.Request, err *AppError) {
	qerror <- fmt.Sprintf("[error] [client %s] %q %s", req.RemoteAddr, req.RequestURI, err.Message)
}

func writeAccess() {
	for x := range qaccess {
		accessLog.Println(x)
	}
}

func writeErrors() {
	for x := range qerror {
		errorLog.Println(x)
	}
}

func logCleanup() {
	close(qaccess)
	close(qerror)
	alf.Close()
	elf.Close()
}

func logReload() {
	a, err := os.OpenFile(AccessLogFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		return
	}
	accessLog = log.New(a, "", log.LstdFlags)
	alf.Close()
	alf = a

	b, err := os.OpenFile(ErrorLogFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		return
	}
	errorLog = log.New(b, "", log.LstdFlags)
	elf.Close()
	elf = b
}
