BIN := "./bin/calendar"
BIN_SENDER := "./bin/sender"
BIN_SCHEDULER := "./bin/scheduler"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/sender
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/scheduler

run: build
	$(BIN) -config-file-path ./configs/config.yaml
	$(BIN_SENDER) -config-file-path ./configs/config_sender.yaml
	$(BIN_SCHEDULER) -config-file-path ./configs/config_scheduler.yaml

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
	go test -race ./...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.63.4

lint: install-lint-deps
	golangci-lint run ./...

migration-create:
	goose -dir migrations create create_event_table sql

migration-up:
	goose -dir migrations postgres 'postgresql://postgres:postgres@localhost:55433/app_db?sslmode=disable' up

migration-down:
	goose -dir migrations postgres 'postgresql://postgres:postgres@localhost:55433/app_db?sslmode=disable' down

migration-status:
	goose -dir migrations postgres 'postgresql://postgres:postgres@localhost:55433/app_db?sslmode=disable' status

generate:
	rm -rf ./internal/server/grpc/pb
	mkdir -p ./internal/server/grpc/pb

	protoc \
		--go_out=./internal/server/grpc \
		--go-grpc_out=./internal/server/grpc \
		./api/*.proto

.PHONY: build run build-img run-img version test lint
