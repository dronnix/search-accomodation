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

PHONY: start-testing ## Create test environment.
start-testing:
	docker-compose -f deployment/docker-compose/test.yaml up -d

PHONY: stop-testing ## Clean up test environment.
stop-testing:
	docker-compose -f deployment/docker-compose/test.yaml down

PHONY: test
test:
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

.PHONY: generate-api
generate-api: ### Generate API structures and servers by swagger spec.
	oapi-codegen -package=api -generate=types -o api/types.go api/geolocation_1.0.0.yaml
	oapi-codegen -package=api -generate chi-server -o api/chi_server.go api/geolocation_1.0.0.yaml

.PHONY: docker-build
docker-build: ### Build docker images.
	bash ./deployment/docker/update-iploc-server.sh
	bash ./deployment/docker/update-iploc-data-importer.sh

.PHONY: compose-run
compose-run: docker-build ### Run apps in docker-compose.
	docker-compose -f deployment/docker-compose/run.yaml up

.PHONY: compose-clean
compose-clean:
	docker-compose -f deployment/docker-compose/run.yaml rm -fs