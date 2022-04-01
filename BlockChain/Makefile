# Labs Makefile
# -------------

.PHONY: run clean test
.PHONY: fmt lint vet codecheck

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test

define USAGE_MSG
Please specify the lab to be build and which assignment part, if it exists.
---------------------------------
Usage: make [OPTION] lab=labX part=partX

---------------------------------
Available OPTIONS:
---------------------------------
* build     - Compile the lab
* test      - Runs the tests
* run       - Runs the lab binary (assume that a main exists)
* codecheck - Runs go format, lint and vet

---------------------------------
Available Labs:
---------------------------------
- lab1
 `- part1
 `- part2
- lab2
 `- part1
 `- part2
- lab3

---------------------------------
Examples:
---------------------------------
  make build lab=lab1 part=part1 
  make test lab=lab2 part=part2  
  make run lab=lab3              
---------------------------------
endef

lab=
part=

ifeq ($(lab),)
 $(info $(USAGE_MSG))
 $(error )
endif

ifeq ($(part),)
LAB_DIR = ./$(lab)
BINARY_PATH=$(LAB_DIR)
else
LAB_DIR = ./$(lab)/$(part)
BINARY_PATH=$(LAB_DIR)/$(part)
endif

build:
	@echo "+ building source for $(lab) $(part)"
	$(GOBUILD) -v -o $(BINARY_PATH) $(LAB_DIR)

# The default is to run all tests of all labs,
# if no lab part is specified
test:
	@echo "+ executing tests for $(lab) $(part)"
	$(GOTEST) -v  $(LAB_DIR)/...

bench:
	@echo "+ executing benchmark for $(lab) $(part)"
	$(GOTEST) -v -run=^$$ -benchmem -bench=. $(LAB_DIR)

run: build
	@echo "+ running $(lab) $(part)"
	$(BINARY_PATH)

clean:
	@echo "+ cleaning $(lab) $(part)"
	@echo $(has_args)
	$(GOCLEAN) -i $(LAB_DIR)/...

codecheck: fmt lint vet

fmt:
	@echo "+ go fmt"
	$(GOCMD) fmt $(LAB_DIR)/...

lint:
	@echo "+ go lint"
	golint -min_confidence=0.1 $(LAB_DIR)/...

vet:
	@echo "+ go vet"
	$(GOCMD) vet $(LAB_DIR)/...
