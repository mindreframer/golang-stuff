Find comprehensive documentation in `/doc/index.html`

A new abstraction for developing and maintaining Big Data applications
======================================================================

The Go Circuit extends the reach of [Go](http://golang.org)'s linguistic
environment to multi-host/multi-process applications.  In simple terms, the Go
Circuit was born from the desire to be able to write:

	feedback := make(chan int)
	circuit.Spawn("host25.datacenter.net", func() {
		feedback <- 1
	})
	<-feedback
	println("Roundtrip complete.")

The `Spawn` operator will start its argument function on a desired
remote host in a new goroutine, while making it possible to communicate between
the parent and child goroutines using the same Go code that you would use to
communicate between traditional goroutines. Above, the channel
`feedback` is transparently “stretched” between the parent
goroutine, executing locally, and the child goroutine, executing remotely and
hosting the anonymous function execution.

Using the circuit one is able to write complex distributed applications —
involving multiple types of collaborating processes — within a single
_circuit program_.  The _circuit language_ used therein is
syntactically identical to Go while also:

* Providing facilities for spawning goroutines on remote hardware, and
* Treating local and remote goroutines in the same manner, both syntactically and semantically.

As a result, distributed application code becomes orders of magnitude shorter,
as compared to using traditional alternatives. For isntance, we have been able
to write large real-world cloud applications — e.g. streaming multi-stage
MapReduce pipelines — in as many as 200 lines of code _from the ground up_.

For lifecycle maintenance, the circuit provides a powerful toolkit that can
introspect into, control and modify various dynamic aspects of a live circuit
application.  Robust networking protocols allow for complex runtime maneuvers
like, for instance, surgically replacing components of running cloud
applications with binaries from different versions of the source tree, without
causing service interruption.

The transparent source of the circuit runtime makes it easy to instrument
circuit deployments with custom logic that has full visibility of cross-runtime
information flow dynamics. Out of the box the circuit comes with a set of tools
for debugging and profiling in-production applications with minimal impact on
uptime.

## License

Copyright 2013 Tumblr, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
