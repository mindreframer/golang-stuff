package gtclient

import (
	"fmt"
	"net"
	"os"
	"time"
)

import (
	l "github.com/ciju/gotunnel/log"
	proto "github.com/ciju/gotunnel/protocol"
	"github.com/ciju/gotunnel/rwtunnel"
)

func ensureServer(addr string) bool {
	lp, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, `
 Local server not running. If your server is,
 running on some other port. Please mention it,
 in the options.

`)
		return false
	}

	lp.Close()
	return true
}

func setupHeartbeat(c net.Conn) {

	for {
		time.Sleep(1 * time.Second)
		c.SetWriteDeadline(time.Now().Add(3 * time.Second))

		_, err := c.Write([]byte("ping"))
		if err != nil {
			l.Log("Couldn't connect to server. Check your network connection, and run client again.")
			os.Exit(1)
		}
	}
}

// connect to server:
// - send the requested subdomain to server.
// - server replies back with a port to setup command channel on.
// - it also replies with the server address that users can access the site on.
func setupCommandChannel(addr, sub string, req, quit chan bool, conn, servInfo chan string) {
	backproxy, err := net.Dial("tcp", addr)
	if err != nil {
		l.Log("CMD: Couldn't connect to ", addr, "err: ", err)
		quit <- true
		return
	}
	defer backproxy.Close()

	proto.SendSubRequest(backproxy, sub)

	// the port to connect on
	serverat, conn_to, _ := proto.ReceiveProxyInfo(backproxy)
	conn <- conn_to
	servInfo <- serverat

	go setupHeartbeat(backproxy)

	for {
		req <- proto.ReceiveConnRequest(backproxy)
	}
}

func SetupClient(port, remote, subdomain string, servInfo chan string) bool {
	localServer := net.JoinHostPort("127.0.0.1", port)

	// if !ensureServer(localServer) {
	// 	return false
	// }

	req, quit, conn := make(chan bool), make(chan bool), make(chan string)

	// fmt.Printf("Setting Gotunnel server %s with local server on %s\n\n", remote, port)

	go setupCommandChannel(remote, subdomain, req, quit, conn, servInfo)

	remoteProxy := <-conn

	// l.Log("remote proxy: %v", remoteProxy)

	for {
		select {
		case <-req:
			// fmt.Printf("New link b/w %s and %s\n", remoteProxy, localServer)
			rp, err := net.Dial("tcp", remoteProxy)
			if err != nil {
				l.Log("Coundn't connect to remote clientproxy", err)
				return false
			}
			lp, err := net.Dial("tcp", localServer)
			if err != nil {
				l.Log("Couldn't connect to localserver", err)
				return false
			}

			go rwtunnel.NewRWTunnel(rp, lp)
		case <-quit:
			return true
		}
	}
	return true
}
