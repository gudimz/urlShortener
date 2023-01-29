# COLORS
green=$(shell echo "\033[32m")
red=$(shell echo "\033[0;31m")
yellow=$(shell echo "\033[0;33m")
end=$(shell echo "\033[0m")

build:
		@echo "$(yellow)building app...$(end)"
		go build -o urlShortener -v ./cmd/api/main.go
		@echo "$(green)built successfully$(end)"
test:
		@echo "$(yellow)start testing shorten...$(end)"
		go test ./internal/shorten/
		@echo "$(green)tested successfully$(end)"

		@echo "$(yellow)start testing server...$(end)"
		go test ./internal/server/
		@echo "$(green)tested successfully$(end)"
docker:
		docker-compose build
run:
		docker-compose up -d
stop:
		docker-compose down
clean:
		rm -rf ./urlShortener

.PHONY: build test docker run stop clean