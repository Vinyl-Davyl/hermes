PREFIX ?= $(HOME)/.local
BIN := bin/hermes

.PHONY: test build install go-install site landing pages-assets clean

test:
	go test ./...

build:
	go build -o $(BIN) ./cmd/hermes

install: build
	mkdir -p $(PREFIX)/bin
	install -m 0755 $(BIN) $(PREFIX)/bin/hermes
	@echo "installed $(PREFIX)/bin/hermes"
	@echo "if hermes is not found, run:  export PATH=\"$(PREFIX)/bin:\$$PATH\""

go-install:
	go install ./cmd/hermes

site: build
	./$(BIN) site

landing: site

pages-assets:
	cp scripts/install.sh web/install
	cp scripts/install.sh web/install.sh

clean:
	rm -rf bin
