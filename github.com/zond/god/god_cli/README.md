cli
===

A simple command line interface to http://github.com/zond/god/client to try out god from a shell.

# Usage

Install with `go get`:

    go get github.com/zond/god/god_cli

Then run from the command line:

    god_cli [-ip 127.0.0.1] [-port 9191] [-enc string] COMMAND

The `-ip` and `-port` options are, not surprisingly, the address and port of a node in the database cluster.

`-enc` is the way the shell string arguments will be encoded to byte slices when sent to the database.

It is one of:

* `string` to simply convert the string to its byte slice representation.
* `float` to convert the string to a big endian 64 bit float in byte slice format.
* `int` to convert the string to a big endian 64 bit int in byte slice format.
* `big` to convert the string to a big endian `math/big.Int`.

If `COMMAND` is ommitted, cli will display the address and position of all nodes in the cluster.

The implemented `COMMAND`s are listed in https://github.com/zond/god/blob/master/god_cli/god_cli.go#L95 and descriptions about them can be found at http://godoc.org/github.com/zond/god/client.
