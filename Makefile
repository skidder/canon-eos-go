GO ?= godep go
COVERAGEDIR = ./coverage
all: clean build test cover

clean: 
	if [ -d $(COVERAGEDIR) ]; then rm -rf $(COVERAGEDIR); fi
	if [ -d bin ]; then rm -rf bin; fi

godep:
	go get github.com/tools/godep

godep-save:
	godep save ./...

all: build test

build:
	if [ ! -d bin ]; then mkdir bin; fi
	$(GO) build ./...

fmt:
	$(GO) fmt ./...

test:
	if [ ! -d $(COVERAGEDIR) ]; then mkdir $(COVERAGEDIR); fi
	$(GO) test -v ./eos -cover -coverprofile=$(COVERAGEDIR)/eos.coverprofile

cover:
	$(GO) tool cover -html=$(COVERAGEDIR)/eos.coverprofile -o $(COVERAGEDIR)/eos.html

bench:
	$(GO) test ./... -cpu 2 -bench .
	