SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

PREFIX := /usr/local
BINDIR := $(PREFIX)/bin
DESTDIR := 

vm:
	go build ./cmd/vm

.PHONY: clean
clean:
	rm vm

.PHONY: install
install: vm
	install -D -m755 -t $(DESTDIR)/$(BINDIR) $^