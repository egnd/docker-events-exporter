#!make

MAKEFLAGS += --always-make

PACKAGES=./internal/... ./pkg/...
BUILD_VERSION=dev

.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

%:
	@:

_ci-conflicts:
	@if grep -rn '^<<<\<<<< ' .; then exit 1; fi
	@if grep -rn '^===\====$$' .; then exit 1; fi
	@if grep -rn '^>>>\>>>> ' .; then exit 1; fi
	@echo "All is OK"

_ci-todos:
	@if grep -rn '@TO\DO:' .; then exit 1; fi
	@echo "All is OK"

_ci-master:
	@git remote update
	@if ! git log --pretty=format:'%H' | grep $$(git log --pretty=format:'%H' -n 1 origin/master) > /dev/null; then exit 1; fi
	@echo "All is OK"

########################################################################################################################

owner: ## Reset folder owner
	sudo chown --changes -R $$(whoami) ./
	@echo "Success"

# mocks: ## Generate mocks
# 	@clear && rm -rf internal/mocks
# 	mockery

# tests: ## Run unit tests
# 	@rm -rf gen/coverage && mkdir -p gen/coverage
# 	CGO_ENABLED=1 go test -mod=readonly -race -cover -covermode=atomic -coverprofile=gen/coverage/profile.out ${PACKAGES} ./cmd/...

# coverage: tests ## Check code coveragem
# 	go tool cover -func=gen/coverage/profile.out
# 	go tool cover -html=gen/coverage/profile.out -o gen/coverage/report.html

# lint: ## Lint source code
# 	@clear
# 	golangci-lint run --color=always --config=.golangci.yml ${PACKAGES}

build: ## Build application
	@mkdir -p bin
	go build -mod=readonly -ldflags "-X 'main.appVersion=$(BUILD_VERSION)-$(GOOS)-$(GOARCH)'" -o bin/app cmd/app/*
	@chmod +x bin/app && ls -lah bin/app && bin/app -version

compose: compose-stop
	docker-compose up --build --abort-on-container-exit --renew-anon-volumes

compose-stop:
	docker-compose down --remove-orphans --volumes

########################################################################################################################

# _env:
# ifeq ($(wildcard .env),)
# 	cp .env.dist .env
# endif

# docker-lint:
# 	docker run --rm -it -v $$(pwd):/src -w /src --entrypoint make golangci/golangci-lint:v1.41 lint

# docker-mocks:
# 	docker run --rm -it -v $$(pwd):/src -w /src --entrypoint sh vektra/mockery:v2 -c "apk add -q make && make mocks"

# docker-tests:
# 	docker run --rm -it -v $$(pwd):/src -w /src --entrypoint make golang:1.17 tests

# docker-coverage:
# 	docker run --rm -it -v $$(pwd):/src -w /src --entrypoint make golang:1.17 coverage
# 	@echo "Read report at file://$$(pwd)/gen/coverage/report.html"

# docker-vendor:
# 	docker run --rm -it -v $$(pwd):/src -w /src --env-file=.env --entrypoint go golang:1.17 mod vendor

# docker-build:
# 	docker build --tag=iptv/schedule:local .
# 	docker run --rm iptv/schedule:local --version
