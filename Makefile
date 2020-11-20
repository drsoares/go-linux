GO    := CGO_ENABLED=0 GOOS=linux GOPROXY=off go
PKG    = github.com/drsoares/go-linux
pkgs   = $(shell $(GO) list $(PKG)/... | grep -v /vendor/)

all: format build

format:
	@echo ">> formatting code"
	@$(GO) fmt $(pkgs)

build:
	@echo ">> building binaries"
	@$(GO) build -mod=vendor $(PKG)/cmd/go-ps

.PHONY: all format build