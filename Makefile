GO    := CGO_ENABLED=0 GOOS=linux go
PKG    = github.com/drsoares/go-ps
pkgs   = $(shell $(GO) list $(PKG)/... | grep -v /vendor/)

all: format build

format:
	@echo ">> formatting code"
	@$(GO) fmt $(pkgs)

build:
	@echo ">> building binaries"
	@$(GO) build $(PKG)/cmd/go-ps

.PHONY: all format build vendorize