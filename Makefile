# Other contants
NAMESPACE=keycloak
PROJECT=keycloak-operator
PKG=github.com/keycloak/keycloak-operator
OPERATOR_SDK_VERSION=v0.10.0
OPERATOR_SDK_DOWNLOAD_URL=https://github.com/operator-framework/operator-sdk/releases/download/$(OPERATOR_SDK_VERSION)/operator-sdk-$(OPERATOR_SDK_VERSION)-x86_64-linux-gnu

# Compile constants
COMPILE_TARGET=./tmp/_output/bin/$(PROJECT)
GOOS=linux
GOARCH=amd64
CGO_ENABLED=0

##############################
# Operator Management        #
##############################
.PHONY: cluster/prepare
cluster/prepare:
	- kubectl apply -f deploy/crds/
	- oc new-project $(NAMESPACE)
	- kubectl apply -f deploy/role.yaml -n $(NAMESPACE)
	- kubectl apply -f deploy/role_binding.yaml -n $(NAMESPACE)
	- kubectl apply -f deploy/service_account.yaml -n $(NAMESPACE)

.PHONY: cluster/clean
cluster/clean:
	# Remove all roles, rolebindings and service accounts with the name keycloak-operator
	- kubectl get roles,rolebindings,serviceaccounts keycloak-operator -n $(NAMESPACE) --no-headers=true -o name | xargs kubectl delete -n $(NAMESPACE)
	# Remove all CRDS with keycloak.org in the name 
	- kubectl get crd --no-headers=true -o name | awk '/keycloak.org/{print $1}' | xargs kubectl delete
	- kubectl delete namespace $(NAMESPACE)

.PHONY: cluster/create/examples
cluster/create/examples:
	- kubectl create -f deploy/examples/keycloak/keycloak.yaml -n $(NAMESPACE)
	- kubectl create -f deploy/examples/realm/realm.yaml -n $(NAMESPACE)

##############################
# Tests                      #
##############################
.PHONY: test/unit
test/unit:
	@echo Running tests:
	go test -v -race -cover -mod=vendor ./pkg/...

##############################
# Local Development          #
##############################
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

.PHONY: code/lint
code/lint:
	@echo "--> Running golangci-lint"
	@which golangci-lint 2>/dev/null ; if [ $$? -eq 1 ]; then \
		go get -u github.com/golangci/golangci-lint/cmd/golangci-lint; \
	fi
	golangci-lint run

##############################
# CI                         #
##############################
.PHONY: setup/travis
setup/travis:
	@echo Installing Operator SDK
	@curl -Lo operator-sdk ${OPERATOR_SDK_DOWNLOAD_URL} && chmod +x operator-sdk && sudo mv operator-sdk /usr/local/bin/

.PHONY: code/coverage
code/coverage:
	@echo "--> Running go coverage"
	@which cover 2>/dev/null ; if [ $$? -eq 1 ]; then \
		go get golang.org/x/tools/cmd/cover; \
	fi
	@go test -coverprofile cover.out -mod=vendor ./pkg/...
	@go tool cover -html=cover.out -o cover.html

.PHONY: code/goveralls
code/goveralls: code/coverage
	go get -u github.com/mattn/goveralls
	@echo "--> Running goveralls"
	@goveralls -coverprofile=cover.out -service=travis-ci
