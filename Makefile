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
	go run main.go generate --package-json package.json

.PHONY: Makefile
Makefile: ## Create a Makefile
	go run main.go generate --makefile Makefile

.PHONY: all
all: ## Run all the scripts
	go run main.go generate --package-json package.json --makefile Makefile

.PHONY: install
install: ## build the binary and copy to ~/.local/bin/
	go build -o ~/.local/bin/metatask main.go

.PHONY: metatask
metatask: ## Build the metatask binary
	go build -o metatask main.go

## metatask-end