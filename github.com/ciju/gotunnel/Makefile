build:
	go build -ldflags "-X main.version `cat VERSION`" -o bin/gotunnel client.go

package: build
	BINS=`ls bin/` && git stash && git checkout binaries && cp bin/* ./ && git add $$BINS && git commit -m"updating binaries" && git checkout master && git stash apply
