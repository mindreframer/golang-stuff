# g + h = github [![Build Status](https://drone.io/github.com/jingweno/gh/status.png)](https://drone.io/github.com/jingweno/gh/latest)

![gh](http://owenou.com/gh/images/gangnamtocat.png)

Fast GitHub command line client. Current version is [0.6.0](https://drone.io/github.com/jingweno/gh/files).

## Overview

gh is a command line client to GitHub. It's designed to run as fast as possible with easy installation across operating systems.
If you like gh, please also take a look at [hub](https://github.com/defunkt/hub). Hub is a reference implementation to gh.

## Motivation

**Fast** 

    $ time hub version > /dev/null
    hub version > /dev/null  0.03s user 0.01s system 93% cpu 0.047 total

    $ time gh version > /dev/null
    gh version > /dev/null  0.01s user 0.01s system 85% cpu 0.022 total

    $ time hub browse > /dev/null
    hub browse > /dev/null  0.07s user 0.04s system 87% cpu 0.130 total

    $ time gh browse > /dev/null
    gh browse > /dev/null  0.03s user 0.02s system 87% cpu 0.059 total

**Muti-platforms**

gh is fully implemented in the Go language and is designed to run across operating systems.

**Easy installation**

There're no pre-requirements to run gh. Download the [binary](https://drone.io/github.com/jingweno/gh/files) and go!

**Unix**

gh commands are single, unhyphenated words that map to their Unix ancestors’ names and flags where applicable.

## Installation

There are [compiled binary forms of gh](https://drone.io/github.com/jingweno/gh/files) for Darwin, Linux and Windows.

To install gh on OSX with [Homebrew](https://github.com/mxcl/homebrew), run:

    $ brew install https://raw.github.com/jingweno/gh/master/homebrew/gh.rb

## Compilation

To compile gh from source, you need to have a [Go development environment](http://golang.org/doc/install), version 1.1 or better, and run:

    $ go get github.com/jingweno/gh

Note that `go get` will pull down sources from various VCS.
Please make sure you have git and hg installed.

## Upgrade

Since gh is under heavy development, I roll out new releases often.
Please take a look at the [CI server](https://drone.io/github.com/jingweno/gh/files) for the latest built binaries.
I plan to implement automatic upgrade in the future.

To upgrade gh on OSX with Homebrew, run:

    $ brew upgrade https://raw.github.com/jingweno/gh/master/homebrew/gh.rb

To upgrade gh from source, run:

    $ go get -u github.com/jingweno/gh

## Usage

### gh help
    
    $ gh help
    Usage: gh [command] [options] [arguments]

    Commands:

        pull              Open a pull request on GitHub
        fork              Make a fork of a remote repository on GitHub and add as remote
        ci                Show CI status of a commit
        browse            Open a GitHub page in the default browser
        compare           Open a compare page on GitHub
        help              Show help
        version           Show gh version

    See 'gh help [command]' for more information about a command.

### gh pull

    # while on a topic branch called "feature":
    $ gh pull
    [ opens text editor to edit title & body for the request ]
    [ opened pull request on GitHub for "YOUR_USER:feature" ]

    # explicit pull base & head:
    $ gh pull -b jingweno:master -h jingweno:feature

    $ gh pull -i 123
    [ attached pull request to issue #123 ]

### gh fork

    $ gh fork
    [ repo forked on GitHub ]
    > git remote add -f YOUR_USER git@github.com:YOUR_USER/CURRENT_REPO.git

    $ gh fork --no-remote
    [ repo forked on GitHub ]

### gh ci

    $ gh ci
    > (prints CI state of HEAD and exits with appropriate code)
    > One of: success (0), error (1), failure (1), pending (2), no
    > status (3)

    $ gh ci BRANCH
    > (prints CI state of BRANCH and exits with appropriate code)
    > One of: success (0), error (1), failure (1), pending (2), no
    > status (3)

    $ gh ci SHA
    > (prints CI state of SHA and exits with appropriate code)
    > One of: success (0), error (1), failure (1), pending (2), no
    > status (3)

### gh browse

    gh browse
    > open https://github.com/YOUR_USER/CURRENT_REPO

    $ gh browse commit/SHA
    > open https://github.com/YOUR_USER/CURRENT_REPO/commit/SHA

    $ gh browse issues
    > open https://github.com/YOUR_USER/CURRENT_REPO/issues

    $ gh browse -u jingweno -r gh
    > open https://github.com/jingweno/gh

    $ gh browse -u jingweno -r gh commit/SHA
    > open https://github.com/jingweno/gh/commit/SHA

    $ git browse -r resque
    > open https://github.com/YOUR_USER/resque

    $ git browse -r resque network
    > open https://github.com/YOUR_USER/resque/network

### gh compare

    $ gh compare refactor
    > open https://github.com/CURRENT_REPO/compare/refactor

    $ gh compare 1.0..1.1
    > open https://github.com/CURRENT_REPO/compare/1.0...1.1

    $ gh compare -u other-user patch
    > open https://github.com/other-user/REPO/compare/patch

## Release Notes

* **0.6.0** June 11, 2013
  * Implement `fork`
* **0.5.2** June 8, 2013
  * Extract GitHub API related code to [`octokat`](https://github.com/jingweno/octokat)
* **0.5.1** June 7, 2013
  * Remove `-p` flag from `browse`
* **0.5.0** June 5, 2013
  * Rename `pull-request` to `pull`
  * Rename `ci-status` to `ci`
* **0.4.1** June 2, 2013
  * Add Rake task to bump version
* **0.4.0** June 2, 2013
  * Implement `compare`
  * Fix bugs on `browse`
* **0.0.3** June 1, 2013
  * Implement `browse`
* **0.0.2** May 29, 2013
  * Implement `ci`
* **0.0.1** May 22, 2013
  * Implement `pull-request`

## Roadmap

* authentication (done)
* gh pull-request (done)
* gh ci-status (done)
* gh browse (done)
* gh compare (done)
* gh fork (done)
* gh clone (in progress)
* gh remote add
* gh fetch
* gh cherry-pick
* gh am, gh apply
* gh check
* gh merge
* gh create
* gh init
* gh push
* gh submodule

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

## License

gh is released under the MIT license. See [LICENSE.md](https://github.com/jingweno/gh/blob/master/LICENSE.md).
