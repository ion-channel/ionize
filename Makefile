# System Setup
SHELL = bash

# Go Stuff
GOCMD=go
GOLINTCMD=golint
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOLIST=$(GOCMD) list
GOVET=$(GOCMD) vet
GOTEST=$(GOCMD) test -v ./...
GOFMT=$(GOCMD) fmt
CGO_ENABLED ?= 0
GOOS ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')

# General Vars
APP := $(shell basename $(PWD) | tr '[:upper:]' '[:lower:]')
DATE := $(shell date -u +%Y-%m-%d%Z%H:%M:%S)
VERSION := v0.1.7
COVERAGE_DIR=coverage

TRAVIS_BUILD_NUMBER ?= 1
TRAVIS_COMMIT ?= $(shell git rev-parse HEAD)

BUILD_NUMBER ?= $(TRAVIS_BUILD_NUMBER)
BUILD_VERSION := $(VERSION)-$(BUILD_NUMBER)
GIT_COMMIT_HASH ?= $(TRAVIS_COMMIT)


.PHONY: all
all: test build

.PHONY: travis_setup
travis_setup: ## Setup the travis environmnet
	@if [[ -n "$$BUILD_ENV" ]] && [[ "$$BUILD_ENV" == "testing" ]]; then echo -e "$(INFO_COLOR)THIS IS EXECUTING AGAINST THE TESTING ENVIRONMEMNT$(NO_COLOR)"; fi
	@sudo /etc/init.d/postgresql stop
	@echo "Installing AWS cli"
	@curl "https://s3.amazonaws.com/aws-cli/awscli-bundle.zip" -o "awscli-bundle.zip"
	@unzip awscli-bundle.zip
	@sudo /usr/bin/python3.4 ./awscli-bundle/install -i /usr/local/aws -b /usr/local/bin/aws
	@echo "Downloading latest Ionize"
	@wget --quiet https://s3.amazonaws.com/public.ionchannel.io/files/ionize/linux/bin/ionize
	@chmod +x ionize && mkdir -p $$HOME/.local/bin && mv ionize $$HOME/.local/bin
	@echo "Installing Go Linter"
	@go get -u golang.org/x/lint/golint

.PHONY: analyze
analyze:  ## Perform an analysis of the project
	@if [[ -n "$$BUILD_ENV" ]] && [[ "$$BUILD_ENV" == "testing" ]]; then \
		IONCHANNEL_SECRET_KEY=$$TESTING_APIKEY IONCHANNEL_ENDPOINT_URL=$$TESTING_ENDPOINT_URL ionize --config .ionize.test.yaml analyze; \
	else \
		ionize analyze; \
	fi

.PHONY: build
build: ## Build the project
	@echo "Building latest image"
	docker build \
		--build-arg BUILD_PATH=/go$${PWD/$$GOPATH} \
		-t ionchannel/ionize .

.PHONY: deploy
deploy: #build
	@echo "Logging into Docker Hub"
	-@echo "$(DOCKER_PASSWORD)" | docker login -u "$(DOCKER_USERNAME)" --password-stdin

	@if [[ ! -z "$(TRAVIS_TAG)" ]] ; then  \
		echo "Pushing release image to Docker Hub" ; \
		docker tag ionchannel/ionize:latest ionchannel/ionize:$(TRAVIS_TAG) ; \
		docker push ionchannel/ionize:$(TRAVIS_TAG) ; \
	fi

	@if [[ "master" == "$(TRAVIS_BRANCH)" ]] ; then \
		echo "Pushing image to Docker Hub" ; \
		docker push ionchannel/ionize:latest ; \
	fi

.PHONY: clean
clean:  ## Clean out all generated files
	-@$(GOCLEAN)
	-@rm -f output.txt
	-@rm -rf deploy
	-@rm -rf coverage
	-@rm -f coverage.txt

.PHONY: coverage
coverage:  ## Generates the code coverage from all the tests
	@echo "Total Coverage: $$(make --no-print-directory coverage_compfriendly | tee coverage.txt)%"

.PHONY: coverage_compfriendly
coverage_compfriendly:  ## Generates the code coverage in a computer friendly manner
	-@rm -rf coverage
	-@mkdir -p $(COVERAGE_DIR)/tmp
	@for j in $$(go list ./... | grep -v '/vendor/' | grep -v '/ext/'); do go test -covermode=count -coverprofile=$(COVERAGE_DIR)/$$(basename $$j).out $$j > /dev/null 2>&1; done
	@echo 'mode: count' > $(COVERAGE_DIR)/tmp/full.out
	@tail -q -n +2 $(COVERAGE_DIR)/*.out >> $(COVERAGE_DIR)/tmp/full.out
	@$(GOCMD) tool cover -func=$(COVERAGE_DIR)/tmp/full.out | tail -n 1 | sed -e 's/^.*statements)[[:space:]]*//' -e 's/%//'

.PHONY: crosscompile
crosscompile:  ## Build the binaries for the primary OS'
	GOOS=linux $(GOBUILD) -ldflags "-X main.buildTime=$(DATE) -X main.appVersion=$(BUILD_VERSION)" -o deploy/linux/bin/$(APP) .
	GOOS=darwin $(GOBUILD) -ldflags "-X main.buildTime=$(DATE) -X main.appVersion=$(BUILD_VERSION)" -o deploy/darwin/bin/$(APP) .
	GOOS=windows $(GOBUILD) -ldflags "-X main.buildTime=$(DATE) -X main.appVersion=$(BUILD_VERSION)" -o deploy/windows/bin/$(APP).exe .

.PHONY: help
help:  ## Show This Help
	@for line in $$(cat Makefile | grep "##" | grep -v "grep" | sed  "s/:.*##/:/g" | sed "s/\ /!/g"); do verb=$$(echo $$line | cut -d ":" -f 1); desc=$$(echo $$line | cut -d ":" -f 2 | sed "s/!/\ /g"); printf "%-30s--%s\n" "$$verb" "$$desc"; done

.PHONY: install
install:  ## Installs the binary
	$(GOCMD) install

.PHONY: test
test: unit_test  ## Runs all available tests

.PHONY: unit_test
unit_test:  ## Run unit tests
	$(GOTEST)

.PHONY: fmt
fmt: ## Run gofmt
	@echo "checking formatting..."
	@$(GOFMT) $(shell $(GOLIST) ./... | grep -v '/vendor/')

.PHONY: vet
vet: ## Run go vet
	@echo "vetting..."
	@$(GOVET) $(shell $(GOLIST) ./... | grep -v '/vendor/')

.PHONY: lint
lint: ## Run golint
	@echo "linting..."
	@$(GOLINTCMD) -set_exit_status $(shell $(GOLIST) ./... | grep -v '/vendor/')
