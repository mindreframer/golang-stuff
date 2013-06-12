TOP=$(shell pwd)
WORK=$(TOP)/work
GO_CHECKOUT=$(WORK)/go
GO_10_ROOT=$(GO_CHECKOUT).10
GO_10_BIN=$(GO_10_ROOT)/bin/go
GO_11_ROOT=$(GO_CHECKOUT).11
GO_11_BIN=$(GO_11_ROOT)/bin/go

# uncomment to benchmark with gccgo
# TESTFLAGS=-compiler=gccgo

BENCHCMP=$(GO_11_ROOT)/misc/benchcmp

# setup our benchmarking environment
GOPATH=$(TOP)
export GOPATH
unexport GOROOT GOBIN

bench: go1 runtime http

go1: $(WORK)/go1-10.txt $(WORK)/go1-11.txt
	$(BENCHCMP) $^

runtime: $(WORK)/runtime-10.txt $(WORK)/runtime-11.txt
	$(BENCHCMP) $^

http: $(WORK)/http-10.txt $(WORK)/http-11.txt
	$(BENCHCMP) $^

update: $(GO_CHECKOUT) $(GO_10_ROOT) $(GO_11_ROOT)
	hg pull --cwd $(GO_CHECKOUT) -u
	hg pull --cwd $(GO_10_ROOT) -u
	hg pull --cwd $(GO_11_ROOT) -u
	rm -rf $(GO_10_ROOT)/bin $(GO_11_ROOT)/bin
	rm -f $(WORK)/*.txt

$(GO_CHECKOUT):
	hg clone https://code.google.com/p/go $@

$(GO_10_ROOT): $(GO_CHECKOUT)
	hg clone -b release-branch.go1 $(GO_CHECKOUT) $@
	hg import --cwd $@ --no-commit $(TOP)/patches/6501099.diff

$(GO_11_ROOT): $(GO_CHECKOUT)
	hg clone -b release-branch.go1.1 $(GO_CHECKOUT) $@

$(GO_10_BIN): $(GO_10_ROOT)
	cd $(GO_10_ROOT)/src ; ./make.bash

$(GO_11_BIN): $(GO_11_ROOT)
	cd $(GO_11_ROOT)/src ; ./make.bash

$(WORK)/go1-10.txt: $(GO_10_BIN)
	$(GO_10_BIN) test $(TESTFLAGS) -bench=. bench/go1 > $@

$(WORK)/go1-11.txt: $(GO_11_BIN)
	$(GO_11_BIN) test $(TESTFLAGS) -bench=. bench/go1 > $@

$(WORK)/runtime-10.txt: $(GO_10_BIN)
	$(GO_10_BIN) test $(TESTFLAGS) -test.run=XXX -test.bench=. bench/runtime > $@

$(WORK)/runtime-11.txt: $(GO_11_BIN)
	$(GO_11_BIN) test $(TESTFLAGS) -test.run=XXX -test.bench=. bench/runtime > $@

$(WORK)/http-10.txt: $(GO_10_BIN)
	$(GO_10_BIN) test $(TESTFLAGS) -test.run=XXX -test.bench=. bench/http > $@

$(WORK)/http-11.txt: $(GO_11_BIN)
	$(GO_11_BIN) test $(TESTFLAGS) -test.run=XXX -test.bench=. bench/http > $@

clean:	
	rm -f $(WORK)/*.txt
