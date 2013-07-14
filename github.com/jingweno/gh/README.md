# g + h = github [![Build Status](https://drone.io/github.com/jingweno/gh/status.png)](https://drone.io/github.com/jingweno/gh/latest)

![gh](http://owenou.com/gh/images/gangnamtocat.png)

Fast GitHub command line client implemented in Go. Current version is [0.15.0](http://bit.ly/go-gh).

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

There're no pre-requirements to run gh. Download the [binary](http://bit.ly/go-gh) and go!

## Installation


### Homebrew

Installing on OSX is easiest with [Homebrew](https://github.com/mxcl/homebrew):

    $ brew install https://raw.github.com/jingweno/gh/master/homebrew/gh.rb

### Standalone

`gh` is easily installed as an executable.
Download the [compiled binary forms of gh](http://bit.ly/go-gh) for Darwin, Linux and Windows.

### Source

To compile gh from source, you need to have a [Go development environment](http://golang.org/doc/install), version 1.1 or better, and run:

    $ go get github.com/jingweno/gh

Note that `go get` will pull down sources from various VCS.
Please make sure you have git and hg installed.

## Upgrade

Since gh is under heavy development, I roll out new releases often.
Please take a look at the [built binaries](http://bit.ly/go-gh) for the latest built binaries.
I plan to implement automatic upgrade in the future.

### Homebrew

To upgrade gh on OSX with Homebrew, run:

    $ brew upgrade https://raw.github.com/jingweno/gh/master/homebrew/gh.rb

### Source

To upgrade gh from source, run:

    $ go get -u github.com/jingweno/gh

## Aliasing

It's best to use `gh` by aliasing it to `git`.
All git commands will still work with `gh` adding some sugar.

`gh alias` displays instructions for the current shell. With the `-s` flag,
it outputs a script suitable for `eval`.

You should place this command in your `.bash_profile` or other startup
script:

    eval "$(gh alias -s)"

For more details, run `gh help alias`.

## Usage

### gh help
    
    $ gh help
    [display help for all commands]
    $ gh help pull-request
    [display help for pull-request]

### gh init

    $ gh init -g
    > git init
    > git remote add origin git@github.com:YOUR_USER/REPO.git

### gh checkout

    $ gh checkout https://github.com/jingweno/gh/pull/35
    > git remote add -f -t feature git://github:com/foo/gh.git
    > git checkout --track -B foo-feature foo/feature

    $ gh checkout https://github.com/jingweno/gh/pull/35 custom-branch-name

### gh merge

    $ gh merge https://github.com/jingweno/gh/pull/73
    > git fetch git://github.com/jingweno/gh.git +refs/heads/feature:refs/remotes/jingweno/feature
    > git merge jingweno/feature --no-ff -m 'Merge pull request #73 from jingweno/feature...'

### gh clone

    $ gh clone jingweno/gh
    > git clone git://github.com/jingweno/gh

    $ gh clone -p jingweno/gh
    > git clone git@github.com:jingweno/gh.git

    $ gh clone jekyll_and_hype
    > git clone git://github.com/YOUR_LOGIN/jekyll_and_hype.

    $ gh clone -p jekyll_and_hype
    > git clone git@github.com:YOUR_LOGIN/jekyll_and_hype.git

### gh fetch

    $ gh fetch jingweno
    > git remote add jingweno git://github.com/jingweno/REPO.git
    > git fetch jingweno

    $ git fetch jingweno,foo
    > git remote add jingweno ...
    > git remote add foo ...
    > git fetch --multiple jingweno foo

    $ git fetch --multiple jingweno foo
    > git remote add jingweno ...
    > git remote add foo ...
    > git fetch --multiple jingweno foo

### gh remote

    $ gh remote add jingweno
    > git remote add -f jingweno git://github.com/jingweno/CURRENT_REPO.git

    $ gh remote add -p jingweno
    > git remote add -f jingweno git@github.com:jingweno/CURRENT_REPO.git

    $ gh remote add origin
    > git remote add -f YOUR_USER git://github.com/YOUR_USER/CURRENT_REPO.git    

### gh pull-request

    # while on a topic branch called "feature":
    $ gh pull-request
    [ opens text editor to edit title & body for the request ]
    [ opened pull request on GitHub for "YOUR_USER:feature" ]

    # explicit pull base & head:
    $ gh pull-request -b jingweno:master -h jingweno:feature

    $ gh pull-request -i 123
    [ attached pull request to issue #123 ]

### gh fork

    $ gh fork
    [ repo forked on GitHub ]
    > git remote add -f YOUR_USER git@github.com:YOUR_USER/CURRENT_REPO.git

    $ gh fork --no-remote
    [ repo forked on GitHub ]

### gh create

    $ gh create
    ... create repo on github ...
    > git remote add -f origin git@github.com:YOUR_USER/CURRENT_REPO.git

    # with description:
    $ gh create -d 'It shall be mine, all mine!'

    $ gh create recipes
    [ repo created on GitHub ]
    > git remote add origin git@github.com:YOUR_USER/recipes.git

    $ gh create sinatra/recipes
    [ repo created in GitHub organization ]
    > git remote add origin git@github.com:sinatra/recipes.git

### gh ci-status

    $ gh ci-status
    > (prints CI state of HEAD and exits with appropriate code)
    > One of: success (0), error (1), failure (1), pending (2), no status (3)

    $ gh ci-status BRANCH
    > (prints CI state of BRANCH and exits with appropriate code)
    > One of: success (0), error (1), failure (1), pending (2), no status (3)

    $ gh ci-status SHA
    > (prints CI state of SHA and exits with appropriate code)
    > One of: success (0), error (1), failure (1), pending (2), no status (3)
    
### gh browse

    $ gh browse
    > open https://github.com/YOUR_USER/CURRENT_REPO

    $ gh browse commit/SHA
    > open https://github.com/YOUR_USER/CURRENT_REPO/commit/SHA

    $ gh browse issues
    > open https://github.com/YOUR_USER/CURRENT_REPO/issues

    $ gh browse -u jingweno -r gh
    > open https://github.com/jingweno/gh

    $ gh browse -u jingweno -r gh commit/SHA
    > open https://github.com/jingweno/gh/commit/SHA

    $ gh browse -r resque
    > open https://github.com/YOUR_USER/resque

    $ gh browse -r resque network
    > open https://github.com/YOUR_USER/resque/network

### gh compare

    $ gh compare refactor
    > open https://github.com/CURRENT_REPO/compare/refactor

    $ gh compare 1.0..1.1
    > open https://github.com/CURRENT_REPO/compare/1.0...1.1

    $ gh compare -u other-user patch
    > open https://github.com/other-user/REPO/compare/patch

## Release Notes

See [Releases](https://github.com/jingweno/gh/releases).

## Roadmap

See [Issues](https://github.com/jingweno/gh/issues?labels=feature&page=1&state=open).

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

## Contributors

See [Contributors](https://github.com/jingweno/gh/graphs/contributors).

## License

gh is released under the MIT license. See [LICENSE.md](https://github.com/jingweno/gh/blob/master/LICENSE.md).
