# golang settings
GO ?= go
GOLINT ?= golangci-lint
GOOS := $(shell $(GO) env GOOS)
GOARCH := $(shell $(GO) env GOARCH)
BIN := $(abspath ./bin/$(GOOS)_$(GOARCH))
GO_ENV ?= GOPRIVATE=github.com/tusmasoma GOBIN=$(BIN)
GOPATH := $(shell go env GOPATH)

# maicroservices
SERVICES := catalog customer order gateway
SERVICE_PATH_PREFIX := go/services

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

.PHONY: test
test:
ifdef SERVICE
	@go test -v -vet=off -race=false -timeout=10m -count=1 ./$(SERVICE_PATH_PREFIX)/$(SERVICE)/...
else
	@go list ./go/... | grep -v /mock | xargs go test -v -vet=off -race=false -timeout=10m -count=1
endif

# golangci-lint: lint for all under the PKG
.PHONY: lint
lint: $(BIN)/golangci-lint-$(GOLANGCI_LINT_VERSION)
ifdef SERVICE
	@echo "Running lint for service: $(SERVICE)"
	cd ./$(SERVICE_PATH_PREFIX)/$(SERVICE) && \
	$(BIN)/golangci-lint run -c ../../../.golangci.yml ./...
else
	@for service in $(SERVICES); do \
		echo "Running lint for service: $$service"; \
		(cd ./$(SERVICE_PATH_PREFIX)/$$service && \
		$(BIN)/golangci-lint run -c ../../../.golangci.yml ./...) || true; \
	done
	@echo "Running lint for pkg/"
	(cd ./pkg && $(BIN)/golangci-lint run -c ../.golangci.yml ./...) || true;
endif

.PHONY: lint-diff
lint-diff: $(BIN)/golangci-lint-$(GOLANGCI_LINT_VERSION)
ifdef SERVICE
	@echo "Running lint-diff for service: $(SERVICE)"
	cd $(SERVICE_PATH_PREFIX)/$(SERVICE) && \
	$(BIN)/golangci-lint run -c ../../../.golangci.yml ./... | reviewdog -f=golangci-lint -diff="git diff origin/main"
else
	@for service in $(SERVICES); do \
		echo "Running lint-diff for service: $$service"; \
		(cd $(SERVICE_PATH_PREFIX)/$$service && \
		$(BIN)/golangci-lint run -c ../../../.golangci.yml ./... | reviewdog -f=golangci-lint -diff="git diff origin/main") || true; \
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
	FILES=$$(find ./go/pkg -type f -name "*.go") && \
	${GO_ENV} $(BIN)/goimports -local "github.com/tusmasoma/pkg" -w $${FILES} && \
	${GO_ENV} $(BIN)/gofumpt -l -w $${FILES};
endif

# proto: generate proto files
.PHONY: proto_gen
proto_gen:
	buf generate

format-proto:
	clang-format -i proto/**/*.proto

.PHONY: generate
generate: generate-deps
	@echo "Running go generate for all packages"
	PATH="$(BIN):$(PATH)" ${GO_ENV} go generate ./...
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

setup: proto_gen generate

build-docs: ${CURDIR}/docs/archive ${CURDIR}/node_modules
build-docs: PROTOC_INCLUDE := -I. -I${GOPATH}/src/github.com/gogo/protobuf -I=${GOPATH}/src/github.com/gogo/protobuf/protobuf
build-docs: build-proto-docs build-openapi
	rm -f docs-error.log
	mkdir -p docs/archive
	npm run build-openapi-web

	@if [ -s "docs-error.log" ]; then\
		echo "failed to build docs error";\
		cat docs-error.log;\
		exit 1;\
	fi

build-openapi:
	cd ./go/services/gateway/cmd/ && swag init -g "../web/web.go" --d "../web" --parseDependency -o "../../../../docs/openapi/web" -ot yaml --v3.1

build-proto-docs:
	buf generate --template '{"version":"v2","plugins":[{"remote":"buf.build/community/pseudomuto-doc:v1.5.1","out":"docs/archive","opt":"html,web-proto.html"}]}' --path proto/web
