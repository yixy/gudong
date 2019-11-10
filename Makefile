.PHONY: build clean

VERSION=v1.0.0
BIN=gudong
EXE=gudong.exe
LDFLAGS_PATH=github.com/yixy/gudong/cmd

GO_ENV=CGO_ENABLED=1
GO_FLAGS=-ldflags="-X $(LDFLAGS_PATH).Ver=$(VERSION) -X '$(LDFLAGS_PATH).Env=`uname -mv`' -X '$(LDFLAGS_PATH).BuildTime=`date`'"
GO=env $(GO_ENV) go
GO_WIN=env GOOS=linux GOARCH=amd64 go

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

build_win:
	mkdir target
	$(GO_WIN) build $(GO_FLAGS) -o target/$(EXE) .