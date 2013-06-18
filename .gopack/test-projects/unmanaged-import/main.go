package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Howdy"))
		})
	http.Handle("/", r)
	http.ListenAndServe(":2345", nil)
}
