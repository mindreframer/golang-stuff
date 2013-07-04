# Hopwatch, a debugging tool for Go

Hopwatch is a simple tool in HTML5 that can help debug Go programs. 
It works by communicating to a WebSockets based client in Javascript.
When your program calls the Break function, it sends debug information to the browser page and waits for user interaction.
Using the functions Display, Printf or Dump (go-spew), you can log information on the browser page.
On the hopwatch page, the developer can view debug information and choose to resume the execution of the program.

[Documentation on godoc.org](http://go.pkgdoc.org/github.com/emicklei/hopwatch)

[![Build Status](https://travis-ci.org/emicklei/hopwatch.png)](https://travis-ci.org/emicklei/hopwatch)

&copy; 2012-2013, http://ernestmicklei.com. MIT License

![hopwath with source](https://s3.amazonaws.com/public.philemonworks.com/hopwatch_with_source.png)