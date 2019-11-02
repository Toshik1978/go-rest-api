GOPATH		= $(shell go env GOPATH)
BUILD_TIME	= $(shell date +"%Y%m%d.%H%M%S")
GIT_VERSION	= $(shell git rev-list -1 HEAD)

.PHONY: all modules prereq mock build lint test test+race test+ci docker clean
.DEFAULT_GOAL := all

all: test build

modules:
	@go mod download

prereq:
	@curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin
	@go get -u github.com/golang/mock/mockgen
	@go get -tags 'postgres' -u github.com/golang-migrate/migrate/v4/cmd/migrate

mock:
	@go generate ./...

build:
	@go build -ldflags "-X main.BuildTime=$(BUILD_TIME) -X main.GitVersion=$(GIT_VERSION)"

lint:
	@golangci-lint run ./... --build-tags=test

test: lint
	@go test ./... -tags=test

test+race: lint
	@go test ./... -race -tags=test

test+ci: lint
	@go test ./... -coverprofile=coverage.txt -covermode=atomic -tags=test

docker:
	@docker build -f Dockerfile.production -t go-rest-api .

clean:
	@go clean
	@rm -f ./go-rest-api
	@rm -f ./coverage.txt
