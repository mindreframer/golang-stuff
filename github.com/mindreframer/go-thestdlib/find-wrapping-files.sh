#!/usr/bin/env bash

find . -name '*.go' -exec grep -R '.\{76,\}' {} \; | awk 'BEGIN { FS = ":" } { print $1 }' | sort | uniq
