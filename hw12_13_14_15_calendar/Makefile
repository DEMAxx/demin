BIN := "./bin/calendar"
MIGRATE_BIN := "./bin/migrate"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.64.6

lint: install-lint-deps
	golangci-lint run ./...

migrate: build-migrate
	$(MIGRATE_BIN) -config ./configs/config.toml

build-migrate:
	go build -v -o $(MIGRATE_BIN) ./cmd/migrate

generate:
	rm -rf events/pb
	mkdir -p events/pb

	protoc \
		--go_out=events/pb \
		--go-grpc_out=events/pb \
		api/*.proto

.PHONY: build run build-img run-img version test lint
