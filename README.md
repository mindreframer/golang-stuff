Go-UrlRouter
============

*Efficient URL routing using a Trie data structure.*

[![Build Status](https://travis-ci.org/ant0ine/go-urlrouter.png?branch=master)](https://travis-ci.org/ant0ine/go-urlrouter)

*Note: This package has been merged into [Go-Json-Rest](https://github.com/ant0ine/go-json-rest), where it has been improved and specialized for the REST API use case.*

This Package implements a URL Router, but instead of using the usual "evaluate all the routes and return the first regexp that matches"
strategy, it uses a Trie data structure to perform the routing. This is more efficient, and scales better for a large number of routes. It supports the :param and \*splat placeholders in the route strings.

Install
-------

This package is "go-gettable", just do:

    go get github.com/ant0ine/go-urlrouter

Example
-------

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


More Examples
-------------

- [Web Server](https://github.com/ant0ine/go-urlrouter/blob/master/examples/webserver/main.go) Demo how to use the router with `net/http`
- [Go-Json-Rest](https://github.com/ant0ine/go-json-rest) A quick and easy way to setup a RESTful JSON API

Documentation
-------------

- [Online Documentation (godoc.org)](http://godoc.org/github.com/ant0ine/go-urlrouter)
- [(Blog Post) Better URL Routing ?] (http://blog.ant0ine.com/typepad/2013/02/better-url-routing-golang-1.html)
- [(Blog Post) Playing with Go1.1beta2, it's much faster!] (http://blog.ant0ine.com/typepad/2013/04/go1-1-beta2-much-faster-.html)


Copyright (c) 2013 Antoine Imbert

[MIT License](https://github.com/ant0ine/go-urlrouter/blob/master/LICENSE)


