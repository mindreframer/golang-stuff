// Copyright (C) 2012 Numerotron Inc.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package bingo

import (
	"flag"
	"fmt"
	"github.com/stathat/spitz"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Env int

const (
	EnvProduction Env = iota
	EnvTesting
	EnvDevel
)

var ContentDir string
var Environment Env
var AfterErrorFunc func(Context, *AppError)

func init() {
	var envFlag string
	flag.StringVar(&envFlag, "e", "production", "environment to run in [devel|testing|production]")
	flag.StringVar(&ContentDir, "d", ".", "content root directory (where server can find 'templates/')")
	flag.Parse()

	switch envFlag {
	case "production":
		Environment = EnvProduction
	case "testing":
		Environment = EnvTesting
	case "devel":
		Environment = EnvDevel
	default:
		log.Fatalf("unknown environment %q", envFlag)
	}

	reloadTemplates := Environment != EnvProduction
	DefaultPool = spitz.New(ContentDir+"/templates", reloadTemplates)
	fmt.Printf("bingo content directory: %q\n", ContentDir)
	fmt.Printf("bingo reloading templates? %v\n", reloadTemplates)

	go trapSignals()
}

func trapSignals() {
	c := make(chan os.Signal, 10)
	signal.Notify(c, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	for {
		x := <-c
		if x == syscall.SIGHUP {
			fmt.Printf("!!! signal %s trapped, reopening log files\n", x)
			logReload()
		} else {
			fmt.Printf("!!! signal %s trapped, exiting\n", x)
			cleanup()
			os.Exit(0)
		}
	}
}

func cleanup() {
	logCleanup()
}

func ListenAndServe(addr string) {
	http.ListenAndServe(addr, nil)
}
