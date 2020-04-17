# Other contants
NAMESPACE=keycloak
PROJECT=keycloak-operator
PKG=github.com/keycloak/keycloak-operator
OPERATOR_SDK_VERSION=v0.15.1
OPERATOR_SDK_DOWNLOAD_URL=https://github.com/operator-framework/operator-sdk/releases/download/$(OPERATOR_SDK_VERSION)/operator-sdk-$(OPERATOR_SDK_VERSION)-x86_64-linux-gnu
MINIKUBE_DOWNLOAD_URL=https://storage.googleapis.com/minikube/releases/v1.4.0/minikube-linux-amd64
KUBECTL_DOWNLOAD_URL=https://storage.googleapis.com/kubernetes-release/release/v1.16.0/bin/linux/amd64/kubectl

# Compile constants
COMPILE_TARGET=./tmp/_output/bin/$(PROJECT)
GOOS=${GOOS:-${GOHOSTOS}}
GOARCH=${GOACH:-${GOHOSTARCH}}
CGO_ENABLED=0

##############################
# Operator Management        #
##############################
.PHONY: cluster/prepare
cluster/prepare:
	@kubectl apply -f deploy/crds/ || true
	@kubectl create namespace $(NAMESPACE) || true
	@which oc 2>/dev/null ; if [ $$? -eq 0 ]; then \
		oc project $(NAMESPACE) || true; \
	fi
	@kubectl apply -f deploy/role.yaml -n $(NAMESPACE) || true
	@kubectl apply -f deploy/role_binding.yaml -n $(NAMESPACE) || true
	@kubectl apply -f deploy/service_account.yaml -n $(NAMESPACE) || true

.PHONY: cluster/clean
cluster/clean:
	# Remove all roles, rolebindings and service accounts with the name keycloak-operator
	@kubectl get roles,rolebindings,serviceaccounts keycloak-operator -n $(NAMESPACE) --no-headers=true -o name | xargs kubectl delete -n $(NAMESPACE)
	# Remove all CRDS with keycloak.org in the name 
	@kubectl get crd --no-headers=true -o name | awk '/keycloak.org/{print $1}' | xargs kubectl delete
	@kubectl delete namespace $(NAMESPACE)

.PHONY: cluster/create/examples
cluster/create/examples:
	@kubectl create -f deploy/examples/keycloak/keycloak.yaml -n $(NAMESPACE)
	@kubectl create -f deploy/examples/realm/basic_realm.yaml -n $(NAMESPACE)

##############################
# Tests                      #
##############################
.PHONY: test/unit
test/unit:
	@echo Running tests:
	@go test -v -tags=unit -coverpkg ./... -coverprofile cover-unit.coverprofile -covermode=count ./pkg/...

.PHONY: test/e2e
test/e2e: cluster/prepare
	@echo Running tests:
	@touch deploy/empty-init.yaml
	# This is not recommended way or running the tests (see https://github.com/operator-framework/operator-sdk/blob/master/doc/test-framework/writing-e2e-tests.md#running-go-test-directly-not-recommended)
	# However, this way we will have a consistent way of running tests on Travis and locally. The downside
	# is that Operator testing harness downloads things manually using `go mod` when executing the tests.
	# Here is a corresponding Operator SDK call:
	# operator-sdk test  local --go-test-flags "-tags=integration -coverpkg ./... -coverprofile cover-e2e.coverprofile -covermode=count" --namespace ${NAMESPACE} --up-local --debug --verbose ./test/e2e
	go test -tags=integration -coverpkg ./... -coverprofile cover-e2e.coverprofile -covermode=count -mod=vendor ./test/e2e/... -root=$(PWD) -kubeconfig=$(HOME)/.kube/config -globalMan deploy/empty-init.yaml -namespacedMan deploy/empty-init.yaml -v -singleNamespace -parallel=1 -localOperator

.PHONY: test/coverage/prepare
test/coverage/prepare:
	@echo Preparing coverage file:
	@echo "mode: count" > cover-all.coverprofile
	@tail -n +2 cover-unit.coverprofile >> cover-all.coverprofile
	@tail -n +2 cover-e2e.coverprofile >> cover-all.coverprofile
	@echo Running test coverage generation:
	@which cover 2>/dev/null ; if [ $$? -eq 1 ]; then \
		go get golang.org/x/tools/cmd/cover; \
	fi
	@go tool cover -html=cover-all.coverprofile -o cover.html

.PHONY: test/coverage
test/coverage: test/coverage/prepare
	@go tool cover -html=cover-all.coverprofile -o cover.html

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

.PHONY: setup/operator-sdk
setup/operator-sdk:
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
	operator-sdk generate crds
	# This is a copy-paste part of `operator-sdk generate openapi` command (suggested by the manual)
	which ./bin/openapi-gen > /dev/null || go build -o ./bin/openapi-gen k8s.io/kube-openapi/cmd/openapi-gen
	./bin/openapi-gen --logtostderr=true -o "" -i ./pkg/apis/keycloak/v1alpha1 -O zz_generated.openapi -p ./pkg/apis/keycloak/v1alpha1 -h ./hack/boilerplate.go.txt -r "-"

.PHONY: code/check
code/check:
	@echo go fmt
	go fmt $$(go list ./... | grep -v /vendor/)

.PHONY: code/fix
code/fix:
	# goimport = gofmt + optimize imports
	@which goimports 2>/dev/null ; if [ $$? -eq 1 ]; then \
		go get golang.org/x/tools/cmd/goimports; \
	fi
	@goimports -w `find . -type f -name '*.go' -not -path "./vendor/*"`

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
	@echo Installing Kubectl
	@curl -Lo kubectl ${KUBECTL_DOWNLOAD_URL} && chmod +x kubectl && sudo mv kubectl /usr/local/bin/
	@echo Installing Minikube
	@curl -Lo minikube ${MINIKUBE_DOWNLOAD_URL} && chmod +x minikube && sudo mv minikube /usr/local/bin/
	@echo Booting Minikube up, see Travis env. variables for more information
	@mkdir -p $HOME/.kube $HOME/.minikube
	@touch $KUBECONFIG
	@sudo minikube start --vm-driver=none --kubernetes-version=v1.16.0
	@sudo chown -R travis: /home/travis/.minikube/

.PHONY: test/goveralls
test/goveralls: test/coverage/prepare
	@echo "Preparing goveralls file"
	go get -u github.com/mattn/goveralls
	@echo "Running goveralls"
	@goveralls -v -coverprofile=cover-all.coverprofile -service=travis-ci
