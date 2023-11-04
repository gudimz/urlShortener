PROJECT_DIR = $(CURDIR)
PROJECT_BIN = $(PROJECT_DIR)/bin
BIN_NAME = url-shortener-api
RUN_TYPE = api

.PHONY: build
build:
	go build -o $(PROJECT_BIN)/$(BIN_NAME) -v ./cmd/$(RUN_TYPE)

.PHONY: test
test:
	go test -v ./...

.PHONY: docker
docker:
	docker-compose build

.PHONY: docker-run
docker-run:
	docker-compose up -d

.PHONY: docker-stop
stop:
	docker-compose down

.PHONY: clean
clean:
	rm -rf $(PROJECT_BIN)/$(BIN_NAME)
