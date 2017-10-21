# System Setup
SHELL = bash

# Go Stuff
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test -v $(shell $(GOCMD) list ./... | grep -v /vendor/)
GOFMT=go fmt
CGO_ENABLED ?= 0
GOOS ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')

# General Vars
APP := $(shell basename $(PWD) | tr '[:upper:]' '[:lower:]')
DATE := $(shell date -u +%Y-%m-%d%Z%H:%M:%S)
VERSION := 0.1.0
COVERAGE_DIR=coverage

TRAVIS_BUILD_NUMBER ?= 1
TRAVIS_COMMIT ?= $(shell git rev-parse HEAD)

BUILD_NUMBER ?= $(TRAVIS_BUILD_NUMBER)
BUILD_VERSION := $(VERSION)-$(BUILD_NUMBER)
GIT_COMMIT_HASH ?= $(TRAVIS_COMMIT)


.PHONY: all
all: test build

.PHONY: build
build: fmt ## Build the project
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) $(GOBUILD) -ldflags "-X main.buildTime=$(DATE) -X main.appVersion=$(BUILD_VERSION)" -o $(APP) .

.PHONY: clean
clean:  ## Clean out all generated files
	-@$(GOCLEAN)
	-@rm $(APP)-linux
	-@rm $(APP)-darwin
	-@rm $(APP)-windows

.PHONY: coverage
coverage:  ## Generates the code coverage from all the tests
	@numbers=0; sum=0; for j in $$(go test -cover $$(go list ./... | grep -v '/vendor/') 2>&1 | sed -e 's/\[no\ test\ files\]/0\.0s\ coverage:\ 0%/g' -e 's/[[:space:]]/\ /g' | tr -d "%" | cut -d ":" -f 2 | cut -d " " -f 2); do ((numbers+=1)) && sum=$$(echo $$sum + $$j | bc); done; avg=$$(echo "$$sum / $$numbers" | bc -l); printf "Total Coverage: %.1f%%\n" $$avg

.PHONY: help
help:  ## Show This Help
	@for line in $$(cat Makefile | grep "##" | grep -v "grep" | sed  "s/:.*##/:/g" | sed "s/\ /!/g"); do verb=$$(echo $$line | cut -d ":" -f 1); desc=$$(echo $$line | cut -d ":" -f 2 | sed "s/!/\ /g"); printf "%-30s--%s\n" "$$verb" "$$desc"; done

.PHONY: install
install:  ## Installs the binary
	$(GOCMD) install

release:  ## Build binaiers for the primary OS'
	GOOS=linux $(GOBUILD) -ldflags "-X main.buildTime=$(DATE) -X main.appVersion=$(VERSION)" -o deploy/linux/$(APP) .
	GOOS=darwin $(GOBUILD) -ldflags "-X main.buildTime=$(DATE) -X main.appVersion=$(VERSION)" -o deploy/darwin/$(APP) .
	GOOS=windows $(GOBUILD) -ldflags "-X main.buildTime=$(DATE) -X main.appVersion=$(VERSION)" -o deploy/windows/$(APP).exe .

.PHONY: test
test: unit_test acceptance_test  ## Runs all available tests

.PHONY: unit_test
unit_test:  ## Run unit tests
	$(GOTEST)

.PHONY: acceptance_test
acceptance_test:  ## Run acceptance tests
	cd ext/acceptance_tests && gem install bundler && bundle install && bundle exec cucumber -t ~@expected_failure

.PHONY: fmt
fmt:  ## Run go fmt
	$(GOFMT)
