BUILD_DIR := bin
BINARIES_DIR := cmd

BINARIES := $$(find $(BINARIES_DIR) -maxdepth 1 \( ! -iname "$(BINARIES_DIR)" \) -type d -exec basename {} \;)
IMPORT_PATH := github.com/dronnix/search-accomodation

.PHONY: build
build: ### Build binaries.
	for bin in $(BINARIES); do \
		go build -mod vendor -o $(BUILD_DIR)/$$bin $(IMPORT_PATH)/$(BINARIES_DIR)/$$bin || exit; \
	done

TEST_SELECTOR := $(if $(sel), -run $(sel))
DIR_SELECTOR := $(if $(dir),"./$(dir)","./...")
FAIL_FAST := $(if $(ff), -failfast)

PHONY: unit-test
unit-test:
	go test $(FAIL_FAST) -v -race -mod=vendor -coverprofile coverage.out $(DIR_SELECTOR) $(TEST_SELECTOR)
	go tool cover -func coverage.out
	go tool cover -html coverage.out -o coverage.html
	rm coverage.out

PHONY: bench
bench:
	go test -bench=. ./...

.PHONY: lint
lint: ### Run golangci-lint. Install it using `make install-tools`.
	golangci-lint run --timeout 3m

.PHONY: install-tools
install-tools: ### Install go develop/test tools.
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.43.0