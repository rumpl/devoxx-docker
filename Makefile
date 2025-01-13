TARGET=./bin/devoxx-container
default: run

.PHONY: build
build: ## Build the application
	@go build -o $(TARGET) .

.PHONY: run
run: build ## Run the application
	@./bin/devoxx-container

clean: ## Clean everything
	rm -fr $(shell dirname $(TARGET))

help: ## Show help
	@echo Available commands:
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
