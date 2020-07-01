# Other contants
NAMESPACE=keycloak
PROJECT=keycloak-operator
PKG=github.com/keycloak/keycloak-operator
OPERATOR_SDK_VERSION=v0.18.2
OPERATOR_SDK_DOWNLOAD_URL=https://github.com/operator-framework/operator-sdk/releases/download/$(OPERATOR_SDK_VERSION)/operator-sdk-$(OPERATOR_SDK_VERSION)-x86_64-linux-gnu
MINIKUBE_DOWNLOAD_URL=https://github.com/kubernetes/minikube/releases/download/v1.9.2/minikube-linux-amd64
KUBECTL_DOWNLOAD_URL=https://storage.googleapis.com/kubernetes-release/release/v1.18.0/bin/linux/amd64/kubectl

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
	@kubectl get all -n $(NAMESPACE) --no-headers=true -o name | xargs kubectl delete -n $(NAMESPACE) || true
	@kubectl get roles,rolebindings,serviceaccounts keycloak-operator -n $(NAMESPACE) --no-headers=true -o name | xargs kubectl delete -n $(NAMESPACE) || true
	@kubectl get pv,pvc -n $(NAMESPACE) --no-headers=true -o name | xargs kubectl delete -n $(NAMESPACE) || true
	# Remove all CRDS with keycloak.org in the name
	@kubectl get crd --no-headers=true -o name | awk '/keycloak.org/{print $1}' | xargs kubectl delete || true
	@kubectl delete namespace $(NAMESPACE) || true

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
test/e2e: setup/operator-sdk
	@echo Running e2e local tests:
	operator-sdk test local --go-test-flags "-tags=integration -coverpkg ./... -coverprofile cover-e2e.coverprofile -covermode=count -timeout 0" --operator-namespace $(NAMESPACE) --up-local --debug --verbose ./test/e2e

.PHONY: test/e2e-latest-image
test/e2e-latest-image:
	@echo Running the latest operator image in the cluster:
	# Doesn't need cluster/prepare as it's done by operator-sdk. Uses a randomly generated namespace (instead of keycloak namespace) to support parallel test runs.
	operator-sdk run local ./test/e2e --go-test-flags "-tags=integration -coverpkg ./... -coverprofile cover-e2e.coverprofile -covermode=count" --debug --verbose

.PHONY: test/e2e-local-image cluster/prepare setup/operator-sdk
test/e2e-local-image: cluster/prepare setup/operator-sdk
	@echo Running e2e tests with a fresh built operator image in the cluster:
	docker build . -t keycloak-operator:test
	@echo Running tests:
	operator-sdk test local --go-test-flags "-tags=integration -coverpkg ./... -coverprofile cover-e2e.coverprofile -covermode=count -timeout 0" --image="keycloak-operator:test" --namespace $(NAMESPACE) --up-local --debug --verbose ./test/e2e

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

.PHONY: setup/mod/verify
setup/mod/verify:
	go mod verify

.PHONY: setup/operator-sdk
setup/operator-sdk:
	@echo Installing Operator SDK
	@curl -Lo operator-sdk ${OPERATOR_SDK_DOWNLOAD_URL} && chmod +x operator-sdk && sudo mv operator-sdk /usr/local/bin/

.PHONY: setup/linter
setup/linter:
	@echo Installing Linter
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.26.0

.PHONY: code/run
code/run:
	@operator-sdk run local --watch-namespace=${NAMESPACE}

.PHONY: code/compile
code/compile:
	@GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=${CGO_ENABLED} go build -o=$(COMPILE_TARGET) -mod=vendor ./cmd/manager

.PHONY: code/gen
code/gen:
	operator-sdk generate k8s
	operator-sdk generate crds --crd-version v1beta1
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
	@$(shell go env GOPATH)/bin/golangci-lint run --timeout 10m

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
	@sudo minikube start --vm-driver=none
	@sudo chown -R travis: /home/travis/.minikube/
	sudo ./hack/modify_etc_hosts.sh "keycloak.local"
	@sudo minikube addons enable ingress

.PHONY: test/goveralls
test/goveralls: test/coverage/prepare
	@echo "Preparing goveralls file"
	go get -u github.com/mattn/goveralls
	@echo "Running goveralls"
	@goveralls -v -coverprofile=cover-all.coverprofile -service=travis-ci
