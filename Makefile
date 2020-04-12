MAKEFLAGS += --silent
PORT := 5000

all: clean compile-all

serve:
	PORT=$(PORT) go run cmd/server/main.go
.PHONY: serve

fetch:
	go run cmd/fetch/main.go
.PHONY: fetch

compile-all: compile-server compile-fetch

compile-server:
	@echo Compiling server
	go build -o out/server cmd/server/main.go
.PHONY: compile-server

compile-fetch:
	@echo Compiling fetch
	go build -o out/fetch cmd/fetch/main.go
.PHONY: compile-fetch

clean:
	@echo Cleaning
	rm -rf out/*

################################################################################
# Cross compilation
################################################################################

cross-compile-all:
	$(MAKE) cross-compile GOOS=linux GOARCH=amd64
	$(MAKE) cross-compile GOOS=windows GOARCH=amd64
	$(MAKE) cross-compile GOOS=darwin GOARCH=amd64
.PHONY: cross-compile-all

cross-compile: ensure-goos ensure-goarch
	$(MAKE) cross-compile-server GOOS=$(GOOS) GOARCH=$(GOARCH)
	$(MAKE) cross-compile-fetch GOOS=$(GOOS) GOARCH=$(GOARCH)
.PHONY: cross-compile

cross-compile-server: ensure-goos ensure-goarch
	@echo Compiling server-$(GOOS)-$(GOARCH)
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o out/server-$(GOOS)-$(GOARCH) cmd/server/main.go
.PHONY: cross-compile-server

cross-compile-fetch: ensure-goos ensure-goarch
	@echo Compiling fetch-$(GOOS)-$(GOARCH)
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o out/fetch-$(GOOS)-$(GOARCH) cmd/fetch/main.go
.PHONY: cross-compile-fetch

################################################################################
# Docker
################################################################################

DOCKER_RUN := docker run -it --rm -v $(shell pwd)/.fetch-cache:/.fetch-cache
DOCKER_BUILD := docker build

docker-run-server:
	$(DOCKER_RUN) -p $(PORT):$(PORT) -e PORT=$(PORT) gamesdb-server
.PHONY: docker-run-server

docker-run-fetch:
	$(DOCKER_RUN) gamesdb-fetch
.PHONY: docker-run-fetch

docker-build-server:
	$(DOCKER_BUILD) . -t gamesdb-server -f Dockerfile.web
.PHONY: docker-build-server

docker-build-fetch:
	$(DOCKER_BUILD) . -t gamesdb-fetch -f Dockerfile.fetch
.PHONY: docker-build-fetch


################################################################################
# Heroku
################################################################################

heroku-all: docker-build-fetch docker-run-fetch heroku-push heroku-release

heroku-push: ensure-app
	@echo Pushing to heroku
	heroku container:push web --recursive --app=$(APP)
.PHONY: heroku-push

heroku-release: ensure-app
	@echo Releasing in heroku
	heroku container:release web --app=$(APP)
.PHONY: heroku-release

################################################################################
# Helpers
################################################################################

ensure-goos:
ifndef GOOS
	$(error GOOS is undefined)
endif

ensure-goarch:
ifndef GOARCH
	$(error GOARCH is undefined)
endif

ensure-app:
ifndef APP
	$(error APP is undefined)
endif