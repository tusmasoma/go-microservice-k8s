# golang settings
GO ?= go
GOLINT ?= golangci-lint
GOOS := $(shell $(GO) env GOOS)
GOARCH := $(shell $(GO) env GOARCH)
BIN := $(abspath ./bin/$(GOOS)_$(GOARCH))
GO_ENV ?= GOPRIVATE=github.com/tusmasoma GOBIN=$(BIN)

# maicroservices
SERVICES := catalog customer order commerce-gateway
SERVICE_PATH_PREFIX := microservice-k8s-demo

# tools
$(shell mkdir -p $(BIN))

GOLANGCI_LINT_VERSION := v1.55.2
$(BIN)/golangci-lint-$(GOLANGCI_LINT_VERSION):
	unlink $(BIN)/golangci-lint || true
	$(GO_ENV) ${GO} install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	mv $(BIN)/golangci-lint $(BIN)/golangci-lint-$(GOLANGCI_LINT_VERSION)
	ln -s $(BIN)/golangci-lint-$(GOLANGCI_LINT_VERSION) $(BIN)/golangci-lint

MOCKGEN_VERSION := 1.6.0
$(BIN)/mockgen-$(MOCKGEN_VERSION):
	unlink $(BIN)/mockgen || true
	$(GO_ENV) ${GO} install github.com/golang/mock/mockgen@v$(MOCKGEN_VERSION)
	mv $(BIN)/mockgen $(BIN)/mockgen-$(MOCKGEN_VERSION)
	ln -s $(BIN)/mockgen-$(MOCKGEN_VERSION) $(BIN)/mockgen

GOIMPORTS_VERSION := v0.19.0
$(BIN)/goimports-$(GOIMPORTS_VERSION):
	unlink $(BIN)/goimports || true
	$(GO_ENV) ${GO} install golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION)
	mv $(BIN)/goimports $(BIN)/goimports-$(GOIMPORTS_VERSION)
	ln -s $(BIN)/goimports-$(GOIMPORTS_VERSION) $(BIN)/goimports

GOFUMPT_VERSION := v0.6.0
$(BIN)/gofumpt-$(GOFUMPT_VERSION):
	unlink $(BIN)/gofumpt || true
	$(GO_ENV) ${GO} install mvdan.cc/gofumpt@$(GOFUMPT_VERSION)
	mv $(BIN)/gofumpt $(BIN)/gofumpt-$(GOFUMPT_VERSION)
	ln -s $(BIN)/gofumpt-$(GOFUMPT_VERSION) $(BIN)/gofumpt

# go: test for all under the PKG
.PHONY: test
test:
ifdef SERVICE
	$(GO) test -v -count=1 ./$(SERVICE_PATH_PREFIX)/$(SERVICE)/...
else
	$(GO) test -v -count=1 $(foreach service,$(SERVICES),./$(SERVICE_PATH_PREFIX)/$(service)/...)
endif

# golangci-lint: lint for all under the PKG
.PHONY: lint
lint: $(BIN)/golangci-lint-$(GOLANGCI_LINT_VERSION)
ifdef SERVICE
	@echo "Running lint for service: $(SERVICE)"
	cd ./$(SERVICE_PATH_PREFIX)/$(SERVICE) && \
	$(BIN)/golangci-lint run -c ./.golangci.yml ./...
else
	@for service in $(SERVICES); do \
		echo "Running lint for service: $$service"; \
		(cd ./$(SERVICE_PATH_PREFIX)/$$service && \
		$(BIN)/golangci-lint run -c ./.golangci.yml ./...) || true; \
	done
endif

.PHONY: lint-diff
lint-diff: $(BIN)/golangci-lint-$(GOLANGCI_LINT_VERSION)
ifdef SERVICE
	@echo "Running lint-diff for service: $(SERVICE)"
	cd $(SERVICE_PATH_PREFIX)/$(SERVICE) && \
	$(BIN)/golangci-lint run -c ./.golangci.yml ./... | reviewdog -f=golangci-lint -diff="git diff origin/main"
else
	@for service in $(SERVICES); do \
		echo "Running lint-diff for service: $$service"; \
		(cd $(SERVICE_PATH_PREFIX)/$$service && \
		$(BIN)/golangci-lint run -c ./.golangci.yml ./... | reviewdog -f=golangci-lint -diff="git diff origin/main") || true; \
	done
endif

.PHONY: fmt
fmt: $(BIN)/goimports-$(GOIMPORTS_VERSION) $(BIN)/gofumpt-$(GOFUMPT_VERSION)
ifdef SERVICE
	@echo "Running fmt for service: $(SERVICE)"
	FILES=$$(find $(SERVICE_PATH_PREFIX)/$(SERVICE) -type f -name "*.go") && \
	LOCAL_PKG="github.com/tusmasoma/microservice-k8s-demo/$(SERVICE)" && \
	${GO_ENV} $(BIN)/goimports -local "$${LOCAL_PKG}" -w $${FILES} && \
	${GO_ENV} $(BIN)/gofumpt -l -w $${FILES}
else
	@for service in $(SERVICES); do \
		echo "Running fmt for service: $$service"; \
		FILES=$$(find $(SERVICE_PATH_PREFIX)/$$service -type f -name "*.go") && \
		LOCAL_PKG="github.com/tusmasoma/microservice-k8s-demo/$$service" && \
		${GO_ENV} $(BIN)/goimports -local "$${LOCAL_PKG}" -w $${FILES} && \
		${GO_ENV} $(BIN)/gofumpt -l -w $${FILES}; \
	done
endif

# .PHONY: generate
# generate: generate-deps
# 	@for dir in $$(find $(if $(SERVICE),$(SERVICE_PATH_PREFIX)/$(SERVICE),$(SERVICE_PATH_PREFIX)) -type d | sed '1,1d' | sed 's@./@@') ; do \
# 		if [ -n "$$(git diff --name-only origin/main "$${dir}")" ]; then \
# 			echo "go generate ./$${dir}/..." && \
# 			(cd "$${dir}" && PATH="$(BIN):$(PATH)" ${GO_ENV} ${GO} generate ./...) || exit 1; \
# 		fi; \
# 	done
# 	$(MAKE) fmt
.PHONY: generate
generate: generate-deps
ifdef SERVICE
	@echo "Running generate for service: $(SERVICE)"
	@for dir in $$(find $(SERVICE_PATH_PREFIX)/$(SERVICE) -type d); do \
		if [ -n "$$(git diff --name-only origin/main -- $$dir)" ]; then \
			echo "go generate $$dir/..." && \
			(cd "$$dir" && PATH="$(BIN):$(PATH)" ${GO_ENV} ${GO} generate ./...) || true; \
		fi; \
	done
else
	@for service in $(SERVICES); do \
		echo "Running generate for service: $$service"; \
		for dir in $$(find $(SERVICE_PATH_PREFIX)/$$service -type d); do \
			if [ -n "$$(git diff --name-only origin/main -- $$dir)" ]; then \
				echo "go generate $$dir/..." && \
				(cd "$$dir" && PATH="$(BIN):$(PATH)" ${GO_ENV} ${GO} generate ./...) || true; \
			fi; \
		done; \
	done
endif
	$(MAKE) fmt

.PHONY: generate-deps
generate-deps: $(BIN)/mockgen-$(MOCKGEN_VERSION)

.PHONY: tidy
tidy:
	$(GO) mod tidy -v $(if $(SERVICE),$(SERVICE_PATH_PREFIX)/$(SERVICE),./$(SERVICE_PATH_PREFIX))

.PHONY: build
build:
	$(GO) build -v $(if $(SERVICE),$(SERVICE_PATH_PREFIX)/$(SERVICE)/...,./$(SERVICE_PATH_PREFIX)/...)

.PHONY: bin-clean
bin-clean:
	$(RM) -r $(if $(SERVICE),$(SERVICE_PATH_PREFIX)/$(SERVICE)/bin,./bin)