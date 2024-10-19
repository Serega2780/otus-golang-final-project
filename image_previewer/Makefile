BIN_CALENDAR := "./bin/image_previewer"
DOCKER_IMG="image_previewer:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN_CALENDAR) -ldflags "$(LDFLAGS)" ./cmd

# run: build
# 	$(BIN_CALENDAR) -config ./configs/calendar_config.yaml

# build-img:
# 	docker build \
# 		--build-arg=LDFLAGS="$(LDFLAGS)" \
# 		-t $(DOCKER_IMG) \
# 		-f build/Dockerfile .
#
# run-img: build-img
# 	docker run $(DOCKER_IMG)
#
# version: build
# 	$(BIN_CALENDAR) version

up:
	docker compose -f docker-compose.yaml up -d

down:
	docker compose -f docker-compose.yaml down

test:
	go test -race -count 100 ./internal/lru ./internal/service

integration_test:
	set -e ;\
	docker compose -f docker-compose-test.yaml build ;\
	test_status_code=0 ;\
	docker compose -f docker-compose-test.yaml run integration_tests go test ./integration || test_status_code=$$? ;\
	docker compose -f docker-compose-test.yaml down ;\
	exit $$test_status_code ;

integration_test-cleanup:
	docker compose -f docker-compose-test.yaml down \
        --rmi local \
		--volumes \
		--remove-orphans \
		--timeout 60; \
  	docker compose rm -f


install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.61.0

lint: install-lint-deps
	golangci-lint run ./...

# run build-img run-img

.PHONY: build  version test lint
