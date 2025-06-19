.DEFAULT_GOAL := all

.PHONY: all deps build install clean userdocs unittest acctest test coverage netdump lint precommit

all: precommit

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

deps:
	go version
	go mod tidy

build:
	$(MAKE) deps
	go fmt ./...
	go build -o ./bin/

install:
	$(MAKE) deps
	go fmt ./...
	go install

clean:
	go version
	go clean -testcache
	rm -rf ./bin
	rm -rf /go/bin/terraform-provider-microsoft365

userdocs:
	go generate
	tfplugindocs validate --provider-name microsoft365


unittest:
	$(MAKE) clean
	$(MAKE) install
	TF_ACC=0 go test -p 16 -timeout 10m -v -cover ./... -run "^TestUnit$(TEST)"

acctest:
	$(MAKE) clean
	$(MAKE) install
ifeq ($(USE_PROXY),1)
	HTTP_PROXY=http://127.0.0.1:8080 HTTPS_PROXY=http://127.0.0.1:8080 TF_ACC=1 go test -p 10 -timeout 300m -v ./... -run "^TestAcc$(TEST)"
else
	TF_ACC=1 go test -p 10 -timeout 300m -v ./... -run "^TestAcc$(TEST)"
endif

test:
	$(MAKE) clean
	$(MAKE) install
	TF_ACC=1 go test -p 10 -timeout 300m -v -cover ./...

coverage:
	$(MAKE) clean
	$(MAKE) install
	@echo "Changed files:"
	@gh pr diff --name-only
	@echo "Running tests"
	TF_ACC=0 go test -p 16 -timeout 10m -v -cover -coverprofile=test-coverage.out ./... -run "^TestUnit$(TEST)"
	@echo "Generating coverage report"
	go tool cover -func=test-coverage.out

netdump:
	mitmdump -p 8080 -w /tmp/mitmproxy.dump

lint:
	golangci-lint --version
	golangci-lint run

precommit:
	$(MAKE) clean
	$(MAKE) build
	$(MAKE) lint
	$(MAKE) unittest
	$(MAKE) userdocs