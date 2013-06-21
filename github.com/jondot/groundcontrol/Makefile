SRC = *.go
PKG = groundcontrol README.md groundcontrol.json.sample web support
VERSION=$(shell cat version.go | perl -n -e'/VERSION = "(.*?)"/ && print $$1')

build: $(SRC)
	go build

package: $(PKG)
	mkdir -p groundcontrol-$(VERSION)
	cp -r $(PKG) groundcontrol-$(VERSION)/
	tar -cvzf groundcontrol-$(VERSION).tar.gz  groundcontrol-$(VERSION)
	rm -rf groundcontrol-$(VERSION)/




.PHONY: clean
