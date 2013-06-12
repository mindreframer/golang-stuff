# Go-WebSocket 

Go-WebSocket is a [Go](http://golang.org/) implementation of the
[WebSocket](http://www.rfc-editor.org/rfc/rfc6455.txt) protocol.

Go-WebSocket passes the server tests in the [Autobahn WebSockets Test
Suite](http://autobahn.ws/testsuite) using the application in the [test
subdirectory](https://github.com/garyburd/go-websocket/tree/master/test).

## Installation

    go get github.com/garyburd/go-websocket/websocket

## Documentation

* [Reference](http://godoc.org/github.com/garyburd/go-websocket/websocket)
* [Chat example](https://github.com/garyburd/go-websocket/tree/master/examples/chat)

## License

Go-WebSocket is available under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).

## Related packages

The [gowebsocket](https://github.com/zhangpeihao/gowebsocket) package wraps the
websocket.Conn type with an implementation of the net.Conn interface. 

## Go-WebSocket compared with other packages

<table>
<tr>
<th></th>
<th><a href="http://godoc.org/github.com/garyburd/go-websocket/websocket">Go-WebSocket</a></th>
<th><a href="http://godoc.org/code.google.com/p/go.net/websocket">go.net</a></th>
</tr>
<tr>
<tr><td>Protocol support</td><td>RFC 6455</td><td>RFC 6455, Hixie 76, Hixie 75</td></tr>
<tr><td>Send pings and receive pongs</td><td>Yes</td><td>No</td></tr>
<tr><td>Send close message</td><td>Yes</td><td>No</td></tr>
<tr><td>Limit size of received message</td><td>Yes</td><td>No</td></tr>
<tr><td>Stream messages</td><td>Yes</td><td>No</td></tr>
<tr><td>Specify IO buffer size</td><td>Yes</td><td>No</td></tr>
</table>
