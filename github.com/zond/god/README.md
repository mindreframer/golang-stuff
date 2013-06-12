god
===

god is a scalable, performant, persistent, in-memory data structure server. It allows massively distributed applications to update and fetch common data in a structured and sorted format.

Its main inspirations are Redis and Chord/DHash. Like Redis it focuses on performance, ease of use, and a small, simple yet powerful feature set, while from the Chord/DHash projects it inherits scalability, redundancy, and transparent failover behaviour.

# Try it out

Install <a href="http://golang.org/doc/install">Go</a>, <a href="http://git-scm.com/downloads">git</a>, <a href="http://mercurial.selenic.com/wiki/Download">Mercurial</a> and <a href="http://gcc.gnu.org/install/">gcc</a>, <code>go get github.com/zond/god/god_server</code>, run <code>god_server</code>, browse to <a href="http://localhost:9192/">http://localhost:9192/</a>.

# Documents

HTML documentation: http://zond.github.com/god/

godoc documentation: http://go.pkgdoc.org/github.com/zond/god

# TODO

* Docs
 * Add illustrations to the usage manual
* Benchmark
 * Consecutively start 1-20 instances on equally powerful machines and benchmark against each size
  * Need 20 machines of equal and constant performance. Is anyone willing to lend me this for few days of benchmarking?
 * Add benchmark results to docs
