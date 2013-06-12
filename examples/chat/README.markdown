# Chat Example

This application shows how to use use the
[Go-WebSocket](https://github.com/garyburd/go-websocket) package and
[jQuery](http://jquery.com) to implement a simple web chat application.

## Running the example

The example requires a working Go development environment. The [Getting
Started](http://golang.org/doc/install) page describes how to install the
development environment.

Once you have Go up and running, you can download, build and run the example
using the following commands.

    $ go get github.com/garyburd/go-websocket/examples/chat
    $ cd `go list -f '{{.Dir}}' github.com/garyburd/go-websocket/examples/chat`
    $ go run *.go

