DEFAULT_GOAL := help

BUILD_FOLDER = dist
CRT_FOLDER = ssl/ca

# Build info
CLIENT_VERSION ?= 0.1.0

.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: proto
proto: ## Generate gRPC protobuf bindings
	./scripts/gen-proto

.PHONY: ssl
ssl: ## Generate SSL certificates for secure communications
	./scripts/gen-ca
	./scripts/issue-crt

.PHONY: keeper ## Build the goph-keeper service
keeper:
	go build -o $(BUILD_FOLDER)/$@ cmd/$@/*.go

.PHONY: keepctl ## Build the goph-keeper client
keepctl:
	./scripts/build-client $(CLIENT_VERSION)

.PHONY: download
download: ## Download go.mod dependencies
	echo Downloading go.mod dependencies
	go mod download

.PHONY: run
run: stop ## Run the project in docker compose
	docker compose -f deployments/docker-compose.yaml up -d --build

.PHONY: stop
stop: ## Stop the running project and destroy containers
	docker compose -f deployments/docker-compose.yaml down

.PHONY: clean
clean: stop
	rm -rf $(BUILD_FOLDER) $(CRT_FOLDER)

.PHONY: install-tools
install-tools: ## Install additional linters and dev tools
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest
	pre-commit install

.PHONY: lint
lint: ## Run linters on the source code
	golangci-lint run
	shellcheck --severity=warning ./scripts/*
	hadolint ./build/docker/Dockerfile

.PHONY: test
test: ## Run unit tests
	@go test -v -race ./... -coverprofile=coverage.out.tmp -covermode atomic
	@cat coverage.out.tmp | grep -v -E "(_mock|.pb).go" > coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out
