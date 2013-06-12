package main

import (
	"fmt"
	"github.com/ant0ine/go-urlrouter"
	"log"
	"net/http"
)

func Hello(w http.ResponseWriter, req *http.Request, params map[string]string) {
	fmt.Fprintf(w, "Hello %s", params["name"])
}

func Bonjour(w http.ResponseWriter, req *http.Request, params map[string]string) {
	fmt.Fprintf(w, "Bonjour %s", params["name"])
}

func main() {

	router := urlrouter.Router{
		Routes: []urlrouter.Route{
			urlrouter.Route{
				PathExp: "/hello/:name",
				Dest:    Hello,
			},
			urlrouter.Route{
				PathExp: "/bonjour/:name",
				Dest:    Bonjour,
			},
		},
	}

	router.Start()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		route, params := router.FindRouteFromURL(r.URL)
		handler := route.Dest.(func(http.ResponseWriter, *http.Request, map[string]string))
		handler(w, r, params)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
