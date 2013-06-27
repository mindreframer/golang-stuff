[![Build Status](https://travis-ci.org/cloudfoundry/gorouter.png)](https://travis-ci.org/cloudfoundry/gorouter)

# gorouter

This repository contains the source code for a Go implementation of the Cloud
Foundry router.

This router is now used on CloudFoundry.com, replacing the old implementation.

## Summary

The original router can be found at cloudfoundry/router. The original router is
backed by nginx, that uses Lua code to connect to a Ruby server that -- based
on the headers of a client's request -- will tell nginx whick backend it should
use. The main limitations in this architecture are that nginx does not support
non-HTTP (e.g. traffic to services) and non-request/response type traffic (e.g.
to support WebSockets), and that it requires a round trip to a Ruby server for
every request.

The Go implementation of the Cloud Foundry router is an attempt in solving
these limitations. First, with full control over every connection to the
router, it can more easily support WebSockets, and other types of traffic (e.g.
via HTTP CONNECT). Second, all logic is contained in a single process,
removing unnecessary latency.

## Getting started

The following instructions may help you get started with gorouter in a
standalone environment.

### Setup

```
git clone https://github.com/cloudfoundry/gorouter.git
cd gorouter
git submodule update --init
./bin/go install router/router
gem install nats
```

### Start

```
# Start NATS server in daemon mode
nats-server -d

# Start gorouter
./bin/router
```

### Usage

When gorouter is used in Cloud Foundry, it receives route updates via NATS
after sending `router.start`.  Routes that haven't responded with
`router.register` in 2 minutes (configurable) are pruned unless connection to
NATS was lost.

Therefore, to maintain an active route, it needs to be updated at least every 2 minutes.
The format of these route updates are as follows:

```json
{
  "host": "127.0.0.1",
  "port": 4567,
  "uris": [
    "my_first_url.vcap.me",
    "my_second_url.vcap.me"
  ],
  "tags": {
    "another_key": "another_value",
    "some_key": "some_value"
  }
}
```

Such a message can be sent to both the `router.register` subject to register
URIs, and to the `router.unregister` subject to unregister URIs, respectively.

```
$ nohup ruby -rsinatra -e 'get("/") { "Hello!" }' &
$ nats-pub 'router.register' '{"host":"127.0.0.1","port":4567,"uris":["my_first_url.vcap.me","my_second_url.vcap.me"],"tags":{"another_key":"another_value","some_key":"some_value"}}'
Published [router.register] : '{"host":"127.0.0.1","port":4567,"uris":["my_first_url.vcap.me","my_second_url.vcap.me"],"tags":{"another_key":"another_value","some_key":"some_value"}}'
$ curl my_first_url.vcap.me:8080
Hello!
```

### Instrumentation

Gorouter provides `/varz` and `/healthz` http endpoints for monitoring.

The `/routes` endpoint returns the entire routing table as JSON. Each route has an associated array of host:port entries.

All of the endpoints require http basic authentication, credentials for which
can be acquired through NATS. The `port`, `user` and password (`pass` is the config attribute) can be explicitly set in the gorouter.yml config
file's `status` section.

```
status:
  port: 8080
  user: some_user
  pass: some_password
```

Example interaction with curl:

```
curl -vvv "http://someuser:somepass@127.0.0.1:8080/routes"
* About to connect() to 127.0.0.1 port 8080 (#0)
*   Trying 127.0.0.1...
* connected
* Connected to 127.0.0.1 (127.0.0.1) port 8080 (#0)
* Server auth using Basic with user 'someuser'
> GET /routes HTTP/1.1
> Authorization: Basic c29tZXVzZXI6c29tZXBhc3M=
> User-Agent: curl/7.24.0 (x86_64-apple-darwin12.0) libcurl/7.24.0 OpenSSL/0.9.8r zlib/1.2.5
> Host: 127.0.0.1:8080
> Accept: */*
> 
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Mon, 25 Mar 2013 20:31:27 GMT
< Transfer-Encoding: chunked
< 
{"0295dd314aaf582f201e655cbd74ade5.cloudfoundry.me":["127.0.0.1:34567"],"03e316d6aa375d1dc1153700da5f1798.cloudfoundry.me":["127.0.0.1:34568"]}
```

## Logs

The router's logging is specified in its YAML configuration file, in a [steno configuration format](http://github.com/cloudfoundry/steno#from-yaml-file).
The meanings of the router's log levels are as follows:

* `fatal` - An error has occurred that makes the current request unservicable.
Examples: the router can't bind to its TCP port, a CF component has published invalid data to the router.
* `warn` - An unexpected state has occurred. Examples: the router tried to publish data that could not be encoded as JSON
* `info`, `debug` - An expected event has occurred. Examples: a new CF component was registered with the router, the router has begun
to prune routes for stale droplets.

## Notes

* 03/05/13: Code is now used on CloudFoundry.com.

* 1/25/13: The code in this repository has not yet been used on CloudFoundry.com.

* 1/25/13: While this implementation can easily support WebSocket
  connections it does not yet.
