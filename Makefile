GOCMD=env go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOTOOL=$(GOCMD) tool
GOGET=$(GOCMD) get

BINARY=gallery-downloader
TESTS=./...
COVERAGE_FILE=coverage.out

.PHONY: all test build coverage clean

all: test build

build:
	$(GOBUILD) -o $(BINARY) -v

test:
	$(GOTEST) -race -v $(TESTS)

coverage:
	$(GOTEST) -coverprofile=$(COVERAGE_FILE) $(TESTS)
	$(GOTOOL) cover -html=$(COVERAGE_FILE)

clean:
	$(GOCLEAN)
	rm -f $(BINARY) $(COVERAGE_FILE)

staticcheck:
	go get -u honnef.co/go/tools/cmd/staticcheck
	go mod tidy
	go run honnef.co/go/tools/cmd/staticcheck ./...
