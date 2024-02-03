SHELL     := /bin/bash
TIMEOUT   := 5s
NAME      := go-http-util
LOCAL     := github.com/AIright/go-http-util
PATHS     := `GO111MODULE=on go list -f '{{.Dir}}' ./...`
VERSION   := `git rev-parse --short HEAD`

PWD := $(PWD)
export PATH := $(PWD)/bin:$(PATH)

.PHONY: lint
lint:
	$(info #Running lint from master..)
	golangci-lint run --new-from-rev=origin/master ./...

# run full lint
.PHONY: lint-full
lint-full:
	$(info #Running lint...)
	@golangci-lint run ./...

.PHONY: deps
deps:
	$(info #Checking deps...)
	@go mod tidy

.PHONY: update
update:
	$(info #Updating deps...)
	@go get -d -u ./...

.PHONY: format
format:
	$(info #Formatting code...)
	@gosimports -local $(LOCAL) -w $(PATHS)
	@gofumpt -w $(PATHS)

# install tools binary: linter, mockgen, etc.
.PHONY: tools
tools:
	@cd tools && go mod tidy && go generate -tags tools

.PHONY: test
test:
	@go test -race -count 1 -timeout $(TIMEOUT) ./...
