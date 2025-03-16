TARGET=./bin/devoxx-container
TARGET_DIR=$(shell dirname $(TARGET))
ROOTFS=rootfs

default: build

$(TARGET):
	@go build -o $(TARGET) .

build: $(TARGET) ## Build the application

.PHONY: run
run: build ## Run the application
	@./bin/devoxx-container

clean:  ## Clean everything
	rm -fr $(ROOTFS)
	rm -fr $(TARGET_DIR)

$(ROOTFS): ## Download and extract the rootfs of the alpine image
	@mkdir -p $(ROOTFS)
	docker export $(shell docker create ubuntu) | tar -C $(ROOTFS) -xvf -

help: ## Show help
	@echo Available commands:
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
