GO=go
GOFLAGS=-race
COV_PROFILE=coverage.txt

.DEFAULT_GOAL=build

fmt:
	$(GO) fmt ./...
.PHONY: fmt

vet: fmt
	$(GO) vet ./...
.PHONY: vet

lint: vet
	golint -set_exit_status=1 ./...
.PHONY: lint

test: lint
	$(GO) clean -testcache
	$(GO) test -v -coverprofile="$(COV_PROFILE)" -covermode="atomic" ./...
.PHONY: test

install:
	$(GO) install ./...
.PHONY: install

build:
	$(GO) build $(GOFLAGS) github.com/mdm-code/scanner/...
.PHONY: build

cover: test
	$(GO) tool cover -html="$(COV_PROFILE)"
.PHONY: cover

clean:
	$(GO) clean github.com/mdm-code/scanner/...
	$(GO) mod tidy
	$(GO) clean -testcache -cache
	rm -f $(COV_PROFILE)
.PHONY: clean
