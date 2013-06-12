package main

import (
	"flag"
	"fmt"
	"github.com/zond/god/common"
	"github.com/zond/god/dhash"
	"runtime"
)

const (
	address = "address"
)

var listenIp = flag.String("listenIp", "127.0.0.1", "IP address to listen at.")
var broadcastIp = flag.String("broadcastIp", "127.0.0.1", "IP address to broadcast to the cluster.")
var port = flag.Int("port", 9191, "Port to listen to for net/rpc connections. The next port will be used for the HTTP service.")
var joinIp = flag.String("joinIp", "", "IP address to join.")
var joinPort = flag.Int("joinPort", 9191, "Port to join.")
var verbose = flag.Bool("verbose", false, "Whether the server should be log verbosely to the console.")
var dir = flag.String("dir", address, "Where to store logfiles and snapshots. Defaults to a directory named after the listening ip/port. The empty string will turn off persistence.")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	if *dir == address {
		*dir = fmt.Sprintf("%v_%v", *broadcastIp, *port)
	}
	s := dhash.NewNodeDir(fmt.Sprintf("%v:%v", *listenIp, *port), fmt.Sprintf("%v:%v", *broadcastIp, *port), *dir)
	if *verbose {
		s.AddChangeListener(func(ring *common.Ring) bool {
			fmt.Println(s.Describe())
			return true
		})
		s.AddMigrateListener(func(dhash *dhash.Node, source, destination []byte) bool {
			fmt.Printf("Migrated from %v to %v\n", common.HexEncode(source), common.HexEncode(destination))
			return true
		})
		s.AddSyncListener(func(source, dest common.Remote, pulled, pushed int) bool {
			fmt.Printf("%v pulled %v and pushed %v keys synchronizing with %v\n", source.Addr, pulled, pushed, dest.Addr)
			return true
		})
		s.AddCleanListener(func(source, dest common.Remote, cleaned, pushed int) bool {
			fmt.Printf("%v cleaned %v and pushed %v keys to %v\n", source.Addr, cleaned, pushed, dest.Addr)
			return true
		})
	}
	s.MustStart()
	if *joinIp != "" {
		s.MustJoin(fmt.Sprintf("%v:%v", *joinIp, *joinPort))
	}

	select {}
}
