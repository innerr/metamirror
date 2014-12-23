export GOPATH := $(shell pwd)

.PHONY: env fast timekeeper all

all: timekeeper test

timekeeper:
	go install -v defynetwork.com/build/...

test:
	go test -v defynetwork.com/...
