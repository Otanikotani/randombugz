PROJECT_NAME := "randombugz"
PKG := "github.com/trilogy-group/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY: all dep build clean

all: build

dep: ## Get the dependencies
	@go mod download

build: dep ## Build the binary file
	@go build -i -o build/${PROJECT_NAME} $(PKG)

clean: ## Remove previous build
	@rm -rf build/
