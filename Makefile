GO_PKG_DIRS  := $(subst $(shell go list -e -m),.,$(shell go list ./... ))

all: clean fmt lint
	go build -ldflags="-s -w" -o checker $(GO_PKG_DIRS)

fmt:
	gofmt -s -w $(GO_PKG_DIRS)

lint:
	golangci-lint run -v $(GO_PKG_DIRS)

clean:
	rm -f checker
