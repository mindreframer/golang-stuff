# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.inc

TARG=gokogiri/xml

GOFILES=\
    attribute.go\
	cdata.go\
	element.go\
	text.go\
	nodeset.go\
	
CGOFILES=\
	document.go\
	fragment.go\
	node.go\
	
CGO_OFILES=\
	helper.o

include $(GOROOT)/src/Make.pkg
