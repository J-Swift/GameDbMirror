MAKEFLAGS += --silent

all: clean compile-all

serve:
	go run cmd/server/main.go
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
	@echo Compiling main-$(GOOS)-$(GOARCH)
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o out/main-$(GOOS)-$(GOARCH) main.go
.PHONY: cross-compile

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
