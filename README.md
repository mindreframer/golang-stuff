golem _v0.4.2_
================================
A lightweight extendable Go WebSocket-framework with [client library](https://github.com/trevex/golem_client). 

License
-------------------------
Golem is available under the  [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html)

Installation
-------------------------
```
go get github.com/trevex/golem
```

Client
-------------------------
A [client](https://github.com/trevex/golem_client) is also available and heavily used in the [examples](https://github.com/trevex/golem_examples).
More information on how the client is used can be found in the [client repository](https://github.com/trevex/golem_client).

Simple Example
-------------------------
Server:
```go
type Hello struct {
	From string `json:"from"`
}
type Answer struct {
	Msg string `json:"msg"`
}
func hello(conn *golem.Connection, data *Hello) {
	conn.Emit("answer", &Answer{"Thanks, "+ data.From + "!"})
}
func main() {
	myrouter := golem.NewRouter()
	myrouter.On("hello", hello)
	http.HandleFunc("/ws", myrouter.Handler())
	http.ListenAndServe(":8080", nil)
}
```
Client:
```javascript
var conn = new golem.Connection("ws://127.0.0.1:8080/ws", true);
conn.on("answer", function(data) {
    console.log("Answer: "+data.msg);
});
conn.on("open", function() {
    conn.emit("hello", { from: "Client" });
});
```
Output in client console would be `Thanks, Client!`.

Documentation
-------------------------
The documentation is provided via [godoc](http://godoc.org/github.com/trevex/golem).

Wiki & Tutorials
-------------------------
More informations and insights can be found on the [wiki page](https://github.com/trevex/golem/wiki) along with a tutorial series to learn how to use golem:
* [Getting started](https://github.com/trevex/golem/wiki/Getting-started)
* [Using rooms](https://github.com/trevex/golem/wiki/Using-rooms)
* [Building a Chat application](https://github.com/trevex/golem/wiki/Building-a-chat-application)
* [Handshake authorisation using Sessions](https://github.com/trevex/golem/wiki/Handshake-authorisation-using-Sessions)
* [Using flash as WebSocket fallback](https://github.com/trevex/golem/wiki/Using-flash-as-WebSocket-fallback)
* [Custom protocol using BSON](https://github.com/trevex/golem/wiki/Custom-protocol-using-BSON)
* [Using an extended connection type](https://github.com/trevex/golem/wiki/Using-an-extended-connection-type)

More Examples
-------------------------
Several examples are available in the [example repository](https://github.com/trevex/golem_examples). To use them simply checkout the
repository and make sure you installed (go get) golem before. A more detailed guide on how
to use them is located in their repository.

History
-------------------------
* _v0.1.0_ 
  * Basic API layout and documentation
* _v0.2.0_ 
  * Evented communication system and routing
  * Basic room implementation (lobbies renamed to rooms for clarity)
* _v0.3.0_ 
  * Protocol extensions through Parsers
  * Room manager for collections of rooms
* _v0.4.0_ 
  * Protocol interchangable
  * Several bugfixes
  * Client up-to-date
* _v0.4.2_
  * Connection type can be extended
  * Close added to connection

Special thanks
-------------------------
* [Gary Burd](http://gary.beagledreams.com/) (for the great WebSocket protocol implementation and insights through his examples)
* [Andrew Gallant](http://burntsushi.net/) (for help on golang-nuts mailing list)
* [Kortschak](https://github.com/kortschak) (for help on golang-nuts mailing list)

TODO
-------------------------
* Verbose and configurable logging
* Testing
