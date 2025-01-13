TARGET=./bin/devoxx-container
TARGET_DIR=$(shell dirname $(TARGET))

default: build

.PHONY: build
build: ## Build the application (default)
	@go build -o $(TARGET) .

clean:  ## Clean everything
	rm -fr $(TARGET_DIR)

help: ## Show help
	@echo Available commands:
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
