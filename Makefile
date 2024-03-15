## metatask-start
.PHONY: format
format: ## Format the code
	go fmt

.PHONY: build
build: ## Build the code
	go build

.PHONY: lint
lint: ## Lint the code
	golangci-lint run

.PHONY: package.json
package.json: ## Create a package.json file
	go run main.go generate

.PHONY: Makefile
Makefile: ## Create a Makefile
	go run main.go generate

.PHONY: metatask
metatask: ## Build the metatask binary
	go build -o metatask main.go

## metatask-end