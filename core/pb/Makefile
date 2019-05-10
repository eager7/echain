# Makefile for protobuf, if protoc: Command not found, then install by: sudo apt-get install protobuf-compiler
default: all

PB = $(wildcard *.proto)
GO = $(PB:.proto=.pb.go)

all: $(GO)

.PHONY: target

%.pb.go: %.proto depends
	protoc -I=. -I=$(GOPATH)/src -I=$(GOPATH)/src/github.com/gogo/protobuf/protobuf --gofast_out=plugins=grpc:. $<

.PHONY: depends
depends:
	if [ ! -d $GOPATH/bin/protoc-gen-gofast ]; then \
		GO111MODULE=off go get github.com/gogo/protobuf/protoc-gen-gofast;\
	fi

.PHONY: clean
clean:
	-rm *.pb.go

