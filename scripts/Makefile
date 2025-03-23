.PHONY: all build run gotool clean help

PROJECT_DIR = ..
BINARY="go_community.exe"

export CGO_ENABLED=0
export GOOS=windows
export GOARCH=amd64

all: gotool build

build:
	cd $(PROJECT_DIR) && go build -o $(BINARY) ./

run:
	cd $(PROJECT_DIR) && go run ./

gotool:
	cd $(PROJECT_DIR) && go fmt ./
	cd $(PROJECT_DIR) && go vet ./

clean:
	if exist $(BINARY) del /Q $(BINARY)

help:
	@echo "make - format Go Code, compile and generate binary files"
	@echo "make build - compile Go Code and generate binary files"
	@echo "make run - run Go Code"
	@echo "make clean - remove binary files"
	@echo "make gotool - run 'fmt' and 'vet'"
