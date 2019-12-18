SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

PREFIX := /usr/local
BINDIR := $(PREFIX)/bin
DATADIR := $(PREFIX)/share
MANDIR := $(DATADIR)/man
DESTDIR := 

.PHONY: all
all: bin data

.PHONY: bin
bin: vm

.PHONY: data
data: vm.fish vm.bash vm.1.gz

vm: *.go cmd/vm/*.go
	go build ./cmd/vm

vm.fish:
	go run ./cmd/vm -- --generate-fish-completion > $@

vm.bash:
	go run ./cmd/vm -- --generate-bash-completion >> $@

vm.1:
	go run ./cmd/vm -- --generate-man-page > $@

vm.1.gz: vm.1
	gzip -k $^

.PHONY: clean
clean:
	-rm vm
	-rm vm.1
	-rm vm.1.gz
	-rm vm.fish
	-rm vm.bash


.PHONY: install install-bin install-man install-data uninstall

install: install-bin install-man install-data

install-bin: vm
	install -D -m755 -t $(DESTDIR)/$(BINDIR) $^

install-man: vm.1.gz
	install -D -m644 -t $(DESTDIR)/$(MANDIR)/man1 $^

install-data: vm.fish vm.bash
	install -D -m644 -t $(DESTDIR)/$(DATADIR)/fish/completions vm.fish
	install -d $(DESTDIR)/$(DATADIR)/bash-completion/completions
	install -m644 -T vm.bash $(DESTDIR)/$(DATADIR)/bash-completion/completions/vm

uninstall:
	rm $(DESTDIR)/$(BINDIR)/vm
	rm $(DESTDIR)/$(MANDIR)/man1/vm.1.gz
	rm $(DESTDIR)/$(DATADIR)/bash-completion/completions/vm
	rm $(DESTDIR)/$(DATADIR)/fish/completions/vm.fish