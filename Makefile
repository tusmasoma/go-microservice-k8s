# golang settings
GO ?= go
GOLINT ?= golangci-lint
GOOS := $(shell $(GO) env GOOS)
GOARCH := $(shell $(GO) env GOARCH)
BIN := $(abspath ./bin/$(GOOS)_$(GOARCH))
GO_ENV ?= GOPRIVATE=github.com/tusmasoma GOBIN=$(BIN)

# maicroservices
SERVICES := catalog customer order gateway
SERVICE_PATH_PREFIX := services

# tools
$(shell mkdir -p $(BIN))

GOLANGCI_LINT_VERSION := v1.63.4
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


PROTOC_VERSION := 24.4
PROTOC_ZIP := protoc-$(PROTOC_VERSION)-linux-x86_64.zip
$(BIN)/protoc-$(PROTOC_VERSION):
	@if ! command -v protoc &> /dev/null; then \
		echo "Installing protoc..."; \
		curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/$(PROTOC_ZIP); \
		unzip -o $(PROTOC_ZIP) -d $(HOME)/.local; \
		rm -f $(PROTOC_ZIP); \
	fi

PROTOC_GEN_GO_VERSION := v1.31.0
$(BIN)/protoc-gen-go-$(PROTOC_GEN_GO_VERSION):
	unlink $(BIN)/protoc-gen-go || true
	$(GO_ENV) ${GO} install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	mv $(BIN)/protoc-gen-go $(BIN)/protoc-gen-go-$(PROTOC_GEN_GO_VERSION)
	ln -s $(BIN)/protoc-gen-go-$(PROTOC_GEN_GO_VERSION) $(BIN)/protoc-gen-go

PROTOC_GEN_GO_GRPC_VERSION := v1.3.0
$(BIN)/protoc-gen-go-grpc-$(PROTOC_GEN_GO_GRPC_VERSION):
	unlink $(BIN)/protoc-gen-go-grpc || true
	$(GO_ENV) ${GO} install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)
	mv $(BIN)/protoc-gen-go-grpc $(BIN)/protoc-gen-go-grpc-$(PROTOC_GEN_GO_GRPC_VERSION)
	ln -s $(BIN)/protoc-gen-go-grpc-$(PROTOC_GEN_GO_GRPC_VERSION) $(BIN)/protoc-gen-go-grpc

proto_tools: $(BIN)/protoc-$(PROTOC_VERSION) $(BIN)/protoc-gen-go-$(PROTOC_GEN_GO_VERSION) $(BIN)/protoc-gen-go-grpc-$(PROTOC_GEN_GO_GRPC_VERSION)

# go: test for all under the PKG
.PHONY: test
test:
ifdef SERVICE
	$(GO) test -v -count=1 ./$(SERVICE_PATH_PREFIX)/$(SERVICE)/...
else
	$(GO) test -v -count=1 $(foreach service,$(SERVICES),./$(SERVICE_PATH_PREFIX)/$(service)/...) ./pkg/...
endif

# golangci-lint: lint for all under the PKG
.PHONY: lint
lint: $(BIN)/golangci-lint-$(GOLANGCI_LINT_VERSION)
ifdef SERVICE
	@echo "Running lint for service: $(SERVICE)"
	cd ./$(SERVICE_PATH_PREFIX)/$(SERVICE) && \
	$(BIN)/golangci-lint run -c ../../.golangci.yml ./...
else
	@for service in $(SERVICES); do \
		echo "Running lint for service: $$service"; \
		(cd ./$(SERVICE_PATH_PREFIX)/$$service && \
		$(BIN)/golangci-lint run -c ../../.golangci.yml ./...) || true; \
	done
	@echo "Running lint for pkg/"
	(cd ./pkg && $(BIN)/golangci-lint run -c ../.golangci.yml ./...) || true;
endif

.PHONY: lint-diff
lint-diff: $(BIN)/golangci-lint-$(GOLANGCI_LINT_VERSION)
ifdef SERVICE
	@echo "Running lint-diff for service: $(SERVICE)"
	cd $(SERVICE_PATH_PREFIX)/$(SERVICE) && \
	$(BIN)/golangci-lint run -c ../../.golangci.yml ./... | reviewdog -f=golangci-lint -diff="git diff origin/main"
else
	@for service in $(SERVICES); do \
		echo "Running lint-diff for service: $$service"; \
		(cd $(SERVICE_PATH_PREFIX)/$$service && \
		$(BIN)/golangci-lint run -c ../../.golangci.yml ./... | reviewdog -f=golangci-lint -diff="git diff origin/main") || true; \
	done
	@echo "Running lint-diff for pkg/"
	(cd ./pkg && $(BIN)/golangci-lint run -c ../.golangci.yml ./... | reviewdog -f=golangci-lint -diff="git diff origin/main") || true;
endif

.PHONY: fmt
fmt: $(BIN)/goimports-$(GOIMPORTS_VERSION) $(BIN)/gofumpt-$(GOFUMPT_VERSION)
ifdef SERVICE
	@echo "Running fmt for service: $(SERVICE)"
	FILES=$$(find $(SERVICE_PATH_PREFIX)/$(SERVICE) -type f -name "*.go") && \
	LOCAL_PKG="github.com/tusmasoma/services/$(SERVICE)" && \
	${GO_ENV} $(BIN)/goimports -local "$${LOCAL_PKG}" -w $${FILES} && \
	${GO_ENV} $(BIN)/gofumpt -l -w $${FILES}
else
	@for service in $(SERVICES); do \
		echo "Running fmt for service: $$service"; \
		FILES=$$(find $(SERVICE_PATH_PREFIX)/$$service -type f -name "*.go") && \
		LOCAL_PKG="github.com/tusmasoma/services/$$service" && \
		${GO_ENV} $(BIN)/goimports -local "$${LOCAL_PKG}" -w $${FILES} && \
		${GO_ENV} $(BIN)/gofumpt -l -w $${FILES}; \
	done
	@echo "Running fmt for pkg/"
	FILES=$$(find ./pkg -type f -name "*.go") && \
	${GO_ENV} $(BIN)/goimports -local "github.com/tusmasoma/pkg" -w $${FILES} && \
	${GO_ENV} $(BIN)/gofumpt -l -w $${FILES};
endif

# proto: generate proto files
.PHONY: proto_gen
proto_gen: proto_tools
	@for service in $(SERVICES); do \
		echo "Running proto_gen for service: $$service"; \
		protoc --proto_path=. \
			--go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative \
			./proto/$$service/*.proto; \
	done

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
	@echo "Running generate for pkg/"
	@for dir in $$(find ./pkg -type d); do \
		if [ -n "$$(git diff --name-only origin/main -- $$dir)" ]; then \
			echo "go generate $$dir/..." && \
			(cd "$$dir" && PATH="$(BIN):$(PATH)" ${GO_ENV} ${GO} generate ./...) || true; \
		fi; \
	done
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
