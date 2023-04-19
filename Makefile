BINARY_NAME := ethereum_parser
BUILD_DIR := ./_bin
DOCKER_IMAGE_NAME := ethereum_parser
DOCKER_CONTAINER_NAME := ethereum_parser_container


build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/ethereum_parser

test:
	go test ./...

run:
	go run ./cmd/ethereum_parser/main.go

clean:
	rm -rf $(BUILD_DIR)

build-docker:
	docker build -t $(DOCKER_IMAGE_NAME) .

run-docker:
	docker run -p 8080:8080 --name $(DOCKER_CONTAINER_NAME) $(DOCKER_IMAGE_NAME)

.PHONY: build test run clean