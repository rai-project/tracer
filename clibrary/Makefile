FIRST_GOPATH              := $(firstword $(subst :, ,$(GOPATH)))
PKGS                      := $(shell go list ./... | grep -v /tests | grep -v /xcpb | grep -v /openpgp)
GOFILES_NOVENDOR          := #$(shell find . -name vendor -prune -o -type f -name '*.go' -not -name '*.pb.go' -print)
GOFILES_BUILD             := #$(shell find . -type f -name '*.go' -not -name '*_test.go')
RAI_TRACER_LIB_VERSION    ?= $(shell cat VERSION)
RAI_TRACER_LIB_OUTPUT     ?= rai_tracer
RAI_TRACER_LIB_REVISION   := $(shell cat COMMIT 2>/dev/null || git rev-parse --short=8 HEAD)
RAI_TRACER_STATIC_LIB_SUFFIX := a
RAI_TRACER_SHARED_LIB_SUFFIX := so
# Support reproducible builds by embedding date according to SOURCE_DATE_EPOCH if present
DATE                      := $(shell date -u -d "@$(SOURCE_DATE_EPOCH)" '+%FT%T%z' 2>/dev/null || date -u '+%FT%T%z')
# BUILDFLAGS                := -buildmode=c-archive -ldflags '-linkmode external -w -extldflags "-v -static"' -gcflags="-B"
BUILDSHAREDFLAGS          := -buildmode=c-shared -ldflags '-linkmode external  -s -w -extldflags "-v"' -gcflags="-B"
BUILDSTATICFLAGS          := -buildmode=c-archive -ldflags '-linkmode external -s -w -extldflags "-v -static"' -gcflags="-B"
TESTFLAGS                 ?=
PWD                       := $(shell pwd)
PREFIX                    ?= $(GOPATH)
BINDIR                    ?= $(PREFIX)/bin
GO                        := go
CGO_CFLAGS                :=
GOOS                      ?= $(shell go version | cut -d' ' -f4 | cut -d'/' -f1)
GOARCH                    ?= $(shell go version | cut -d' ' -f4 | cut -d'/' -f2)

build: $(RAI_TRACER_LIB_OUTPUT)

sysinfo:
	@echo ">> SYSTEM INFORMATION"
	@echo -n "     PLATFORM: $(shell uname -a)"
	@printf '%s\n' '$(OK)'
	@echo -n "     PWD:    : $(shell pwd)"
	@printf '%s\n' '$(OK)'
	@echo -n "     GO      : $(shell go version)"
	@printf '%s\n' '$(OK)'
	@echo -n "     BUILDSTATICFLAGS: $(BUILDSTATICFLAGS)"
	@echo -n "     BUILDSHAREDFLAGS: $(BUILDSHAREDFLAGS)"
	@printf '%s\n' '$(OK)'
	@echo -n "     GIT     : $(shell git version)"
	@printf '%s\n' '$(OK)'

deps:
	go get github.com/karalabe/xgo

bootstrap:
	# @GOOS=linux GOARCH=386 CGO_ENABLED=1 CC=486-linux-musl-gcc go install std
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-linux-musl-cc go install std
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go install std
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=i686-w64-mingw32-gcc go install std

setup:
	mkdir -p LibraryResources
	mkdir -p dist/Linux-x86-64/
	mkdir -p dist/MacOSX-x86-64/
	mkdir -p dist/Windows-x86-64/

$(RAI_TRACER_LIB_OUTPUT)-darwin: $(GOFILES_NOVENDOR)
	@echo -n ">> CROSSCOMPILE darwin GOOS=darwin GOARCH=amd64 $(CGO_CFLAGS) $(GO) build -v -o dist/MacOSX-x86-64/lib$(RAI_TRACER_LIB_OUTPUT).$(RAI_TRACER_STATIC_LIB_SUFFIX) $(BUILDSTATICFLAGS) $^"
	@GOOS=darwin GOARCH=amd64 $(CGO_CFLAGS) $(GO) build -o dist/MacOSX-x86-64/lib$(RAI_TRACER_LIB_OUTPUT).$(RAI_TRACER_STATIC_LIB_SUFFIX) $(BUILDSTATICFLAGS) $^

	@echo -n ">> CROSSCOMPILE darwin GOOS=darwin GOARCH=amd64 $(CGO_CFLAGS) $(GO) build -v -o dist/MacOSX-x86-64/lib$(RAI_TRACER_LIB_OUTPUT).$(RAI_TRACER_SHARED_LIB_SUFFIX) $(BUILDSHAREDFLAGS) $^"
	@GOOS=darwin GOARCH=amd64 $(CGO_CFLAGS) $(GO) build -o dist/MacOSX-x86-64/lib$(RAI_TRACER_LIB_OUTPUT).$(RAI_TRACER_SHARED_LIB_SUFFIX) $(BUILDSHAREDFLAGS) $^
	@printf '%s\n' '$(OK)'

$(RAI_TRACER_LIB_OUTPUT)-linux: $(GOFILES_NOVENDOR)
	@echo -n ">> CROSSCOMPILE linux GOOS=linux GOARCH=amd64 $(CGO_CFLAGS) $(GO) build -v -o dist/Linux-x86-64/lib$(RAI_TRACER_LIB_OUTPUT).$(RAI_TRACER_STATIC_LIB_SUFFIX) $(BUILDSTATICFLAGS) $^"
	@GOOS=linux GOARCH=amd64 $(CGO_CFLAGS) $(GO) build -o dist/Linux-x86-64/lib$(RAI_TRACER_LIB_OUTPUT).$(RAI_TRACER_STATIC_LIB_SUFFIX) $(BUILDSTATICFLAGS) $^
	# @echo -n ">> CROSSCOMPILE linux GOOS=linux GOARCH=amd64 $(CGO_CFLAGS) $(GO) build -v -o dist/Linux-x86-64/lib$(RAI_TRACER_LIB_OUTPUT).$(RAI_TRACER_SHARED_LIB_SUFFIX) $(BUILDSHAREDFLAGS) $^"
	# @GOOS=linux GOARCH=amd64 $(GO) build -o dist/Linux-x86-64/lib$(RAI_TRACER_LIB_OUTPUT).$(RAI_TRACER_SHARED_LIB_SUFFIX) $(BUILDSHAREDFLAGS) $^
	@printf '%s\n' '$(OK)'

$(RAI_TRACER_LIB_OUTPUT)-linux-cupti: $(GOFILES_NOVENDOR)
	@echo -n ">> CROSSCOMPILE linux GOOS=linux GOARCH=amd64 $(CGO_CFLAGS) $(GO) build -tags cupti -v -o dist/Linux-x86-64/lib$(RAI_TRACER_LIB_OUTPUT)-cupti.$(RAI_TRACER_STATIC_LIB_SUFFIX) $(BUILDSTATICFLAGS) $^"
	@GOOS=linux GOARCH=amd64 $(CGO_CFLAGS) $(GO) build -tags cupti -o dist/Linux-x86-64/lib$(RAI_TRACER_LIB_OUTPUT)-cupti.$(RAI_TRACER_STATIC_LIB_SUFFIX) $(BUILDSTATICFLAGS) $^
	# @echo -n ">> CROSSCOMPILE linux GOOS=linux GOARCH=amd64 $(CGO_CFLAGS) $(GO) build -v -o dist/Linux-x86-64/lib$(RAI_TRACER_LIB_OUTPUT).$(RAI_TRACER_SHARED_LIB_SUFFIX) $(BUILDSHAREDFLAGS) $^"
	# @GOOS=linux GOARCH=amd64 $(GO) build -o dist/Linux-x86-64/lib$(RAI_TRACER_LIB_OUTPUT).$(RAI_TRACER_SHARED_LIB_SUFFIX) $(BUILDSHAREDFLAGS) $^
	@printf '%s\n' '$(OK)'

$(RAI_TRACER_LIB_OUTPUT)-windows: $(GOFILES_NOVENDOR)
	@echo -n ">> CROSSCOMPILE windows"
	@GOOS=windows GOARCH=amd64 CC=i686-w64-mingw32-gcc $(CGO_CFLAGS) $(GO) build -o dist/Windows-x86-64/$@.dll $(BUILDFLAGS)
	@printf '%s\n' '$(OK)'
	@echo -n ">> CROSSCOMPILE linux/386"
	@GOOS=windows GOARCH=386 CC=i686-w64-mingw32-gcc $(CGO_CFLAGS) $(GO) build -o dist/Windows/$@.dll $(BUILDFLAGS)
	@printf '%s\n' '$(OK)'

$(RAI_TRACER_LIB_OUTPUT): $(RAI_TRACER_LIB_OUTPUT)-$(GOOS)


all: setup $(RAI_TRACER_LIB_OUTPUT)

fmt:
	go fmt ./...

.PHONY: all
