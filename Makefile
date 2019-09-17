# Image build contants
REG=quay.io
ORG=keycloak
PROJECT=keycloak-operator
TAG?=latest

#Compile constants
COMPILE_TARGET=./tmp/_output/bin/$(PROJECT)
GOOS=linux
GOARCH=amd64
CGO_ENABLED=0

#Other contants
NAMESPACE=keycloak
PKG=github.com/keycloak/keycloak-operator
OPERATOR_SDK_VERSION=v0.10.0
OPERATOR_SDK_DOWNLOAD_URL=https://github.com/operator-framework/operator-sdk/releases/download/$(OPERATOR_SDK_VERSION)/operator-sdk-$(OPERATOR_SDK_VERSION)-x86_64-linux-gnu

.PHONY: setup
setup: setup/mod setup/githooks code/gen

.PHONY: setup/githooks
setup/githooks:
	@echo Setting up Git hooks:
	ln -sf $$PWD/.githooks/* $$PWD/.git/hooks/

.PHONY: setup/mod
setup/mod:
	@echo Adding vendor directory
	go mod vendor
	@echo setup complete

.PHONY: setup/travis
setup/travis:
	@echo Installing Operator SDK
	@curl -Lo operator-sdk ${OPERATOR_SDK_DOWNLOAD_URL} && chmod +x operator-sdk && sudo mv operator-sdk /usr/local/bin/

.PHONY: code/run
code/run:
	@operator-sdk up local --namespace=${NAMESPACE}

.PHONY: code/compile
code/compile:
	@GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=${CGO_ENABLED} go build -o=$(COMPILE_TARGET) -mod=vendor ./cmd/manager

.PHONY: code/gen
code/gen:
	operator-sdk generate k8s
	operator-sdk generate openapi

.PHONY: code/check
code/check:
	@echo go fmt
	go fmt $$(go list ./... | grep -v /vendor/)

.PHONY: code/fix
code/fix:
	@gofmt -w `find . -type f -name '*.go' -not -path "./vendor/*"`

.PHONY: image/build/push
image/build/push: image/build image/push

.PHONY: test/unit
test/unit:
	@echo Running tests:
	go test -v -race -cover -mod=vendor ./pkg/...

.PHONY: code/lint
code/lint:
	@echo "--> Running golangci-lint"
	@which golangci-lint 2>/dev/null ; if [ $$? -eq 1 ]; then \
		go get -u github.com/golangci/golangci-lint/cmd/golangci-lint; \
	fi
	golangci-lint run
