# gopack

The [natural logarithm](https://en.wikipedia.org/wiki/Natural_logarithm) to Go's [e](http://en.wikipedia.org/wiki/E_(mathematical_constant). Simple package management for Go a la [rebar](https://github.com/basho/rebar).

A configuration file tells gopack about your dependencies and which version should be included. You can point to a tag, a branch, or, if you are being naughty, master. The programming community would thank you not to carry out such a travesty as it leaves your code open to breaking changes. Much better to point at _immutable_ code.

```toml
[deps.memcache]
import = "github.com/bradfitz/gomemcache/memcache"
tag = "1.2"

[deps.mux]
import = "github.com/gorilla/mux"
branch = "1.0rc2"

[deps.toml]
import = "github.com/pelletier/go-toml"
commit = "23d36c08ab90f4957ae8e7d781907c368f5454dd"
```

Then simply run, install, and test your code much as you would have with the ```go``` command. Just replace ```go``` with ```gp```.

```gp test```
```gp run *.go```

etc…

The ```gp``` command will make sure your dependencies are downloaded, their respective git repos are pointed at the appropriate tag or branch, and your code is compiled against the desired library versions. Project dependencies are stored locally in the ```vendor``` directory.

# Installation

First checkout and build from source
```
git clone git@github.com:d2fn/gopack.git
cd gopack
go get github.com/pelletier/go-toml
go build
```

Then copy the ```gopack``` binary to your project directory and invoke just as you would go. Make sure the current directory is on your path or place the ```gp``` binary elsewhere on your path.
```
cp gopack ~/projects/mygoproject/gp
cd ~/projects/myproject
gp run *.go
```

# License

Copyright (c) 2013 Dietrich Featherston

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
