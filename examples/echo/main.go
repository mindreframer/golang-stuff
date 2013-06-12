package main

import (
	"github.com/garyburd/go-websocket/websocket"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func echo(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r.Header, nil, 1024, 1024)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	defer conn.Close()
	for {
		op, r, err := conn.NextReader()
		if err != nil {
			return
		}
		if op != websocket.OpBinary && op != websocket.OpText {
			continue
		}
		w, err := conn.NextWriter(op)
		if err != nil {
			return
		}
		io.Copy(w, r)
		w.Close()
	}
}

func main() {
	http.HandleFunc("/", echo)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
