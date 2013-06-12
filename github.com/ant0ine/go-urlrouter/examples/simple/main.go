package main

import (
	"fmt"
	"github.com/ant0ine/go-urlrouter"
)

func main() {

	router := urlrouter.Router{
		Routes: []urlrouter.Route{
			urlrouter.Route{
				PathExp: "/resources/:id",
				Dest:    "one_resource",
			},
			urlrouter.Route{
				PathExp: "/resources",
				Dest:    "all_resources",
			},
		},
	}

	err := router.Start()
	if err != nil {
		panic(err)
	}

	input := "http://example.org/resources/123"
	route, params, err := router.FindRoute(input)
	if err != nil {
		panic(err)
	}
	fmt.Print(route.Dest)  // one_resource
	fmt.Print(params["id"])  // 123
}
