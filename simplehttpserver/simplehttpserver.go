package simplehttpserver

import (
	"fmt"
	"net/http"
)

import (
	l "github.com/ciju/gotunnel/log"
)

func NewSimpleHTTPServer(port string, dir string) {

	if dir == "" {
		l.Fatal("No directory given, to serve")
	}

	http.Handle("/", http.FileServer(http.Dir(dir)))

	fmt.Println("Serving", dir, "at port", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		l.Fatal("error", err)
	}
}
