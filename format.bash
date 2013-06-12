#!/usr/bin/env bash
find . -iname '*.go' -print0 | xargs -0 -I file gofmt -tabs=false -tabwidth 4 -w file
