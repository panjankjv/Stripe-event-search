.PHONY: init lint build build-all clean test

BINARY_NAME=stripe-event-search

GO111MODULE=on
VERSION := $(shell git tag --points-at HEAD --sort=-v:refname | head -n 1)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' \
           -X 'main.revision=$(REVISION)'


init:
	go mod download

lint:
	@type golangci-lint > /dev/null || go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint $(LINT_OPT) run ./...

# build binary
build:
	go build -o bin/${BINARY_NAME} ./cmd

build-macos:
	@make _build BUILD_OS=darwin BUILD_ARCH=amd64

build-macos-m1:
	@make _build BUILD_OS=darwin BUILD_ARCH=arm64

build-linux:
	@make _build BUILD_OS=linux BUILD_ARCH=amd64

build-linux-arm:
	@make _build BUILD_OS=linux BUILD_ARCH=arm64

build-windows:
	@make _build BUILD_OS=windows BUILD_ARCH=amd64

_build:
	@mkdir -p bin/release
	$(eval BUILD_OUTPUT := ${BINARY_NAME}_${BUILD_OS}_${BUILD_ARCH}${BUILD_ARM})
	GOOS=${BUILD_OS} \
	GOARCH=${BUILD_ARCH} \
	GOARM=${BUILD_ARM} \
	go build -o bin/${BUILD_OUTPUT} ./cmd
	@if [ "${USE_ARCHIVE}" = "1" ]; then \
		gzip -k -f bin/${BUILD_OUTPUT} ;\
		mv bin/${BUILD_OUTPUT}.gz bin/release/ ;\
	fi

build-all: clean
	@make build-macos build-macos-m1 build-linux build-linux-arm build-windows USE_ARCHIVE=1

clean:
	rm -f bin/${BINARY_NAME}_*
	rm -f bin/release/*

test:
	go test -v ./...
