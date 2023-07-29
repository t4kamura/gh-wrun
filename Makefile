BIN := $(shell go list .)
VERSION := $$(make -s show-version)
LDFLAGS := "-s -w"
GOBIN ?= $(shell go env GOPATH)/bin

.PHONY: build
build:
	go build -ldflags=$(LDFLAGS) -trimpath -o $(BIN) .

.PHONY: xbuild
xbuild: $(GOBIN)/goxz
	goxz -n $(BIN) -pv=v$(VERSION) -arch=amd64 -os linux,darwin,windows -build-ldflags=$(LDFLAGS) .

.PHONY: install
install:
	go install -ldflags=$(LDFLAGS) -trimpath .

.PHONY: show-version
show-version: $(GOBIN)/gobump
	@gobump show -r .

.PHONY: clean
clean:
	rm -rf $(BIN) goxz
	go clean

$(GOBIN)/gobump:
	go install github.com/x-motemen/gobump/cmd/gobump@latest

$(GOBIN)/goxz:
	go install github.com/Songmu/goxz/cmd/goxz@latest
