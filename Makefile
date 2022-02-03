PACKAGE_NAME := "$(shell head -n 1 go.mod | cut -d ' ' -f2)"
BINARY_NAME := wisdom

GOOS=$(go env GOOS)
GOPATH=$(go env GOPATH)

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: all test build

all: help

bin: bin-server bin-client

## Go:
dep: ## Tidy up mod file
	@go mod tidy

## Build server:
bin-server: ## Build the project
	@mkdir -p bin
	@CGO_ENABLED=0 GOOS=${GOOS} go build -o bin/"$(BINARY_NAME)"-server cmd/server/main.go

## Build client:
bin-client: ## Build the project
	@mkdir -p bin
	@CGO_ENABLED=0 GOOS=${GOOS} go build -o bin/"$(BINARY_NAME)"-client cmd/client/main.go

clean: ## Remove build related file
	rm -fr ./bin
	rm -f ./coverage.out

## Test:
test: ## Run the tests of the project
	@go test -tags="unit" -v -race ./...

generate: ## Generate mocks for testing
	@rm -rf ./internal/mocks/*
	@PATH=${GOPATH}/bin:${PATH}
	@go install github.com/golang/mock/mockgen/...@v1.6.0
	@go generate -x ./internal/...

coverage: ## Run the tests of the project with coverage
	@go test -tags="unit" -coverprofile=coverage.out -covermode=count ./...
	@go tool cover -func=coverage.out

lint: ## Lint with revive
	@revive -config ./revive.toml ./...
## Local:
local: ## Run a local instance of the application
	@docker-compose -f docker-compose.local.yaml up --build -d --remove-orphans

## Help:
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)
