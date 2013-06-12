#!/bin/sh

PACKAGE="github.com/eblume/proto"

godoc ${PACKAGE} > documentation.txt
godoc --html -src ${PACKAGE} > documentation.html