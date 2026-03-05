BINARY_SERVER = server
BINARY_CLI = stuff
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS = -ldflags "-s -w -X main.version=$(VERSION)"
GOFLAGS = -trimpath

PLATFORMS = linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

.PHONY: build test clean release deploy install-man install-completion

build:
	go build $(GOFLAGS) $(LDFLAGS) -o $(BINARY_SERVER) ./cmd/server
	go build $(GOFLAGS) $(LDFLAGS) -o $(BINARY_CLI) ./cmd/stuff

test:
	go test ./...

clean:
	rm -f $(BINARY_SERVER) $(BINARY_CLI)
	rm -rf dist/

release: clean test
	@mkdir -p dist
	@for platform in $(PLATFORMS); do \
		OS=$${platform%/*}; \
		ARCH=$${platform#*/}; \
		EXT=""; \
		if [ "$$OS" = "windows" ]; then EXT=".exe"; fi; \
		echo "Building $$OS/$$ARCH..."; \
		GOOS=$$OS GOARCH=$$ARCH go build $(GOFLAGS) $(LDFLAGS) -o dist/$(BINARY_SERVER)-$$OS-$$ARCH$$EXT ./cmd/server; \
		GOOS=$$OS GOARCH=$$ARCH go build $(GOFLAGS) $(LDFLAGS) -o dist/$(BINARY_CLI)-$$OS-$$ARCH$$EXT ./cmd/stuff; \
		if [ -f stuff.1 ]; then \
			tar czf dist/allmystuff-$$OS-$$ARCH.tar.gz \
				-C dist $(BINARY_SERVER)-$$OS-$$ARCH$$EXT $(BINARY_CLI)-$$OS-$$ARCH$$EXT \
				-C .. stuff.1 completions/ README.md LICENSE; \
		else \
			tar czf dist/allmystuff-$$OS-$$ARCH.tar.gz \
				-C dist $(BINARY_SERVER)-$$OS-$$ARCH$$EXT $(BINARY_CLI)-$$OS-$$ARCH$$EXT \
				-C .. README.md LICENSE; \
		fi; \
	done

deploy:
	go install $(GOFLAGS) $(LDFLAGS) ./cmd/stuff

install-man:
	install -d $(DESTDIR)/usr/local/share/man/man1
	install -m 644 stuff.1 $(DESTDIR)/usr/local/share/man/man1/stuff.1

install-completion:
	install -d $(DESTDIR)/usr/local/share/bash-completion/completions
	install -m 644 completions/stuff.bash $(DESTDIR)/usr/local/share/bash-completion/completions/stuff
	install -d $(DESTDIR)/usr/local/share/zsh/site-functions
	install -m 644 completions/stuff.zsh $(DESTDIR)/usr/local/share/zsh/site-functions/_stuff
