.PHONY: build clean

VERSION=0.0.1
BIN=gudong
LDFLAGS_PATH=github.com/yixy/gudong/cmd

GO_ENV=CGO_ENABLED=1
GO_FLAGS=-ldflags="-X $(LDFLAGS_PATH).Ver=$(VERSION) -X '$(LDFLAGS_PATH).Env=`uname -mv`' -X '$(LDFLAGS_PATH).BuildTime=`date`'"
GO=env $(GO_ENV) go

UNAME := $(shell uname)

# build cli
build:
	mkdir target
	$(GO) build $(GO_FLAGS) -o target/$(BIN) .
# test
test:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...
# clean all build result
clean:
	go clean ./...
	rm -rf target
