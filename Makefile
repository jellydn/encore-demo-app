.PHONY: dev
dev:
	encore run

.PHONY: test
test:
	encore test ./...

.PHONY: help
help:
	encore help

.PHONY: install
install:
	curl -L https://encore.dev/install.sh | bash
