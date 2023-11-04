PROJECT_DIR = $(CURDIR)
PROJECT_BIN = $(PROJECT_DIR)/bin
BIN_NAME = url-shortener-api
RUN_TYPE = api

GOLANGCI_TAG = 1.55.2
GOLANGCI_LINT_BIN = $(PROJECT_BIN)/golangci-lint

.PHONY: .install-linter
.install-linter:
	@if [ ! -f $(GOLANGCI_LINT_BIN) ]; then \
		$(info "Downloading golangci-lint v$(GOLANGCI_TAG)") \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(PROJECT_BIN) v$(GOLANGCI_TAG); \
	fi


.PHONY: lint
lint: .install-linter
	$(GOLANGCI_LINT_BIN) run ./... --config=./configs/.golangci.yml

.PHONY: lint-fast
lint-fast: .install-linter
	$(GOLANGCI_LINT_BIN) run ./... --fast --config=./configs/.golangci.yml

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
