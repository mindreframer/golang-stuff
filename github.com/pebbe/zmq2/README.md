A Go interface to [ZeroMQ](http://www.zeromq.org/) version 2.

Requires ZeroMQ version 2.1 or 2.2

For ZeroMQ version 3, see: http://github.com/pebbe/zmq3

Including all examples of [ØMQ - The Guide](http://zguide.zeromq.org/page:all).

## Install

    go get github.com/pebbe/zmq2

## Docs

 * [package help](http://godoc.org/github.com/pebbe/zmq2)

## Support for ZeroMQ version 2.1

 * The following functions are not supported in ZeroMQ version 2.1, and will return an error:
  * (*Socket) GetRcvtimeo
  * (*Socket) GetSndtimeo
  * (*Socket) SetRcvtimeo
  * (*Socket) SetSndtimeo
