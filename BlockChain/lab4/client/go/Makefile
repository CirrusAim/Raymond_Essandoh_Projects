PKG := .

.PHONY: run clean deps
.PHONY: fmt lint vet codecheck

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOBIN_DIR = $(PKG)
GO_BINDINGS_DIR =  $(PKG)/go-bindings
CONTRACTS_SRC = ../../contracts
SOLC_BUILD_DIR = ../../build
CONTRACTS_DEPS = ../../node_modules/@openzeppelin/contracts/,
SOLC_ARGS = @openzeppelin/contracts=../../node_modules/@openzeppelin/contracts

BINARY_FILE = $(GOBIN_DIR)/betting

all: generate build

build:
	@echo "+ building source"
	$(GOBUILD) -v -o $(BINARY_FILE)

test: generate
	@echo "+ executing tests"
	$(GOTEST) -v  $(PKG)/...

run: build
	@echo "+ running"
	$(BINARY_FILE)

clean:
	@echo "+ cleaning"
	$(GOCLEAN) -i $(PKG)/...
	rm -rf $(BINARY_FILE) $(GO_BINDINGS_DIR)

deps:
	@echo "+ installing dependencies"
	$(GOGET) github.com/ethereum/go-ethereum
	$(GOGET) golang.org/x/lint/golint

generate:
	@echo "+ go generate"
	mkdir -p $(GO_BINDINGS_DIR)/betting
	$(GOCMD) generate $(PKG)/...

solc:
	solc $(SOLC_ARGS) --allow-paths $(CONTRACTS_DEPS) --overwrite --combined-json abi,asm,ast,bin,bin-runtime,compact-format,devdoc,hashes,interface,metadata,opcodes,srcmap,srcmap-runtime,userdoc $(CONTRACTS_SRC)/Betting.sol -o $(SOLC_BUILD_DIR)

codecheck: fmt lint vet

fmt:
	@echo "+ go fmt"
	$(GOCMD) fmt $(PKG)/...

lint:
	@echo "+ go lint"
	golint -min_confidence=0.1 $(PKG)/...

vet:
	@echo "+ go vet"
	$(GOCMD) vet $(PKG)/...

# Cross compilation
build-linux:
	@echo "+ building $@"
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GOBUILD) -v -o $(BINARY_FILE)

build-osx: generate
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_FILE) -v

build-windows: generate
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_FILE) -v