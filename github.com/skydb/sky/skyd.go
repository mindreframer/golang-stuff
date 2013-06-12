package main

import (
	"./skyd"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
)

//------------------------------------------------------------------------------
//
// Constants
//
//------------------------------------------------------------------------------

const (
	defaultPort = 8585
	defaultDataDir = "/var/lib/sky"
)

const (
	portUsage = "the port to listen on"
	dataDirUsage = "the data directory"
)

const (
	pidPath = "/var/run/skyd.pid"
)

//------------------------------------------------------------------------------
//
// Variables
//
//------------------------------------------------------------------------------

var port uint
var dataDir string

//------------------------------------------------------------------------------
//
// Functions
//
//------------------------------------------------------------------------------

//--------------------------------------
// Initialization
//--------------------------------------

func init() {
	flag.UintVar(&port, "port", defaultPort, portUsage)
	flag.UintVar(&port, "p", defaultPort, portUsage+"(shorthand)")
	flag.StringVar(&dataDir, "data-dir", defaultDataDir, dataDirUsage)
	flag.StringVar(&dataDir, "d", defaultDataDir, dataDirUsage+"(shorthand)")
}

//--------------------------------------
// Main
//--------------------------------------

func main() {
	// Parse the command line arguments.
	flag.Parse()
	
	// Hardcore parallelism right here.
	runtime.GOMAXPROCS(runtime.NumCPU())
	
	// Initialize
	server := skyd.NewServer(port, dataDir)
	writePidFile()
	//setupSignalHandlers(server)
	
	// Start the server up!
	c := make(chan bool)
	err := server.ListenAndServe(c)
	if err != nil {
		fmt.Printf("%v\n", err)
		cleanup(server)
		return
	}
	<- c
	cleanup(server)
}

//--------------------------------------
// Signals
//--------------------------------------

// Handles signals received from the OS.
func setupSignalHandlers(server *skyd.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
	    for _ = range c {
			fmt.Fprintln(os.Stderr, "Shutting down...")
			cleanup(server)
			fmt.Fprintln(os.Stderr, "Shutdown complete.")
			os.Exit(1)
	    }
	}()
}

//--------------------------------------
// Utility
//--------------------------------------

// Shuts down the server socket and closes the database.
func cleanup(server *skyd.Server) {
	if server != nil {
		server.Shutdown()
	}
	deletePidFile()
}

// Writes a file to /var/run that contains the current process id.
func writePidFile() {
	pid := fmt.Sprintf("%d", os.Getpid())
	if err := ioutil.WriteFile(pidPath, []byte(pid), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to write pid file: %v\n", err)
	}
}

// Deletes the pid file.
func deletePidFile() {
	if _, err := os.Stat(pidPath); !os.IsNotExist(err) {
		if err = os.Remove(pidPath); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to remove pid file: %v\n", err)
		}
	}
}