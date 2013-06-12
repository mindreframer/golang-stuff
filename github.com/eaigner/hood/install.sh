#!/bin/sh
go run cmd/gen/templates.go
go build -o $GOPATH/bin/hood github.com/eaigner/hood/cmd
