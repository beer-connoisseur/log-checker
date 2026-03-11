LOCAL_BIN := $(CURDIR)/bin
GOLANGCI_BIN := $(LOCAL_BIN)/golangci-lint
CUSTOM_GCL_BIN := $(LOCAL_BIN)/custom-gcl
LOGCHECKER_BIN := $(LOCAL_BIN)/logchecker

.PHONY: build
build:
	go build -o $(LOCAL_BIN)/logchecker ./cmd/logchecker

bin/golangci-lint:
	curl -sSfL https://golangci-lint.run/install.sh | sh -s v2.10.1

bin/custom-gcl: bin/golangci-lint
	$(GOLANGCI_BIN) custom --destination ./bin

.PHONY: lint
lint: bin/custom-gcl
	@echo 'Running linter on files...'
	$(CUSTOM_GCL_BIN) run \
	--config=.golangci.yaml \
	--max-issues-per-linter=0 \
	--max-same-issues=0

.PHONY: fix
fix: build
	$(LOGCHECKER_BIN) -fix ./...

