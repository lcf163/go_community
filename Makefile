.PHONY: all build run gotool clean help

BINARY="go-community.exe"

export CGO_ENABLED=0
export GOOS=windows
export GOARCH=amd64

all: gotool build

build:
	go build -o ${BINARY}

run:
	go run ./

gotool:
	go fmt ./
	go vet ./

clean:
	if exist $(BINARY) del /Q $(BINARY)

help:
	@echo "make - format Go Code, compile and generate binary files"
	@echo "make build - compile Go Code and generate binary files"
	@echo "make run - run Go Code"
	@echo "make clean - remove binary files"
	@echo "make gotool - run 'fmt' and 'vet'"
