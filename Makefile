# Other contants
NAMESPACE=keycloak
PROJECT=keycloak-operator
PKG=github.com/keycloak/keycloak-operator
OPERATOR_SDK_VERSION=v0.18.2
ifeq ($(shell uname),Darwin)
  OPERATOR_SDK_ARCHITECTURE=x86_64-apple-darwin
else
  OPERATOR_SDK_ARCHITECTURE=x86_64-linux-gnu
endif
OPERATOR_SDK_DOWNLOAD_URL=https://github.com/operator-framework/operator-sdk/releases/download/$(OPERATOR_SDK_VERSION)/operator-sdk-$(OPERATOR_SDK_VERSION)-$(OPERATOR_SDK_ARCHITECTURE)

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
	@kubectl delete -f deploy/service_account.yaml -n $(NAMESPACE) || true
	@kubectl delete -f deploy/role_binding.yaml -n $(NAMESPACE) || true
	@kubectl delete -f deploy/role.yaml -n $(NAMESPACE) || true
	@kubectl delete namespace $(NAMESPACE) || true
	@kubectl delete -f deploy/crds/ || true
	
.PHONY: cluster/clean/monitoring
cluster/clean/monitoring:
	@kubectl delete -n $(NAMESPACE) --all blackboxtargets
	@kubectl delete -n $(NAMESPACE) --all grafanadashboards
	@kubectl delete -n $(NAMESPACE) --all grafanadatasources
	@kubectl delete -n $(NAMESPACE) --all applicationmonitorings
	@kubectl delete crd grafanas.integreatly.org
	@kubectl delete crd grafanadashboards.integreatly.org
	@kubectl delete crd grafanadatasources.integreatly.org
	@kubectl delete crd blackboxtargets.applicationmonitoring.integreatly.org
	@kubectl delete crd applicationmonitorings.applicationmonitoring.integreatly.org
	@kubectl delete namespace application-monitoring

.PHONY: cluster/prepare/monitoring
cluster/prepare/monitoring:
	oc label namespace $(NAMESPACE) "monitoring-key=middleware"
	$(eval _OS_PROMETHEUS_USER=$(shell oc get secrets -n openshift-monitoring grafana-datasources -o 'go-template={{index .data "prometheus.yaml"}}' | base64 --decode | jq -r '.datasources[0].basicAuthUser'))
	$(eval _OS_PROMETHEUS_PASS=$(shell oc get secrets -n openshift-monitoring grafana-datasources -o 'go-template={{index .data "prometheus.yaml"}}' | base64 --decode | jq -r '.datasources[0].basicAuthPassword'))
	kubectl label namespace $(NAMESPACE) monitoring-key=middleware || true
	git clone --depth=1  git@github.com:integr8ly/application-monitoring-operator.git /tmp/keycloak-operator || true
	$(MAKE) -C /tmp/keycloak-operator cluster/install
	cat ./deploy/examples/monitoring/federation.yaml | sed -e 's/<user>/'"$(_OS_PROMETHEUS_USER)"'/g' | \
		sed -e 's@<pass>@'"$(_OS_PROMETHEUS_PASS)"'@g' > /tmp/keycloak-operator/integreatly-additional.yaml || true
	kubectl create secret generic integreatly-additional-scrape-configs --from-file=/tmp/keycloak-operator/integreatly-additional.yaml --dry-run=client -o yaml | kubectl apply -n application-monitoring -f -
	rm -rf /tmp/keycloak-operator/

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
	@go test -v -tags=unit -coverpkg ./... -coverprofile cover-unit.coverprofile -covermode=count -mod=vendor ./pkg/...

.PHONY: test/e2e
test/e2e: setup/operator-sdk
	@echo Running e2e local tests:
	operator-sdk test local --go-test-flags "-tags=integration -coverpkg ./... -coverprofile cover-e2e.coverprofile -covermode=count -timeout 0" --operator-namespace $(NAMESPACE) --up-local --debug --verbose ./test/e2e

.PHONY: test/e2e-latest-image
test/e2e-latest-image:
	@echo Running the latest operator image in the cluster:
	# Doesn't need cluster/prepare as it's done by operator-sdk. Uses a randomly generated namespace (instead of keycloak namespace) to support parallel test runs.
	operator-sdk run local ./test/e2e --go-test-flags "-tags=integration -coverpkg ./... -coverprofile cover-e2e.coverprofile -covermode=count" --debug --verbose

.PHONY: test/ibm-validation
test/ibm-validation:
	@echo Running the operator image in the cluster
	operator-sdk test local ./test/e2e --go-test-flags "-tags=integration -coverpkg ./... -coverprofile cover-e2e.coverprofile -covermode=count -timeout 0" --operator-namespace $(NAMESPACE) --debug --verbose --global-manifest=deploy/empty-init.yaml --namespaced-manifest=deploy/operator.yaml

.PHONY: test/e2e-local-image setup/operator-sdk
test/e2e-local-image: setup/operator-sdk
	@echo Backing up operator.yaml
	@cp deploy/operator.yaml deploy/operator.yaml_bckp
	@echo Building operator image:
	eval $$(minikube -p minikube docker-env); \
	docker build . -t keycloak-operator:test
	@echo Modifying operator.yaml
	@sed -i 's/imagePullPolicy: Always/imagePullPolicy: Never/g' deploy/operator.yaml
	@echo Creating namespace
	kubectl create namespace $(NAMESPACE) || true
	@echo Running e2e tests with a fresh built operator image in the cluster:
	trap 'mv -f deploy/operator.yaml_bckp deploy/operator.yaml' EXIT; \
	operator-sdk test local --go-test-flags "-tags=integration -coverpkg ./... -coverprofile cover-e2e.coverprofile -covermode=count -timeout 0" --image="keycloak-operator:test" --debug --verbose --operator-namespace $(NAMESPACE) ./test/e2e

.PHONY: test/coverage/prepare
test/coverage/prepare:
	@echo Preparing coverage file:
	@echo "mode: count" > cover-all.coverprofile
	@echo "mode: count" > cover-e2e.coverprofile
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
code/gen: client/gen
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

.PHONY: client/gen
client/gen:
	@echo "--> Running code-generator to generate clients"
	# prepare tool code-generator
	@mkdir -p ./tmp/code-generator
	@git clone https://github.com/kubernetes/code-generator.git --branch v0.21.0-alpha.2 --single-branch  ./tmp/code-generator
	# generate client
	./tmp/code-generator/generate-groups.sh "client,informer,lister" github.com/keycloak/keycloak-operator/pkg/client github.com/keycloak/keycloak-operator/pkg/apis keycloak:v1alpha1 --output-base ./tmp --go-header-file ./hack/boilerplate.go.txt
	# check generated client at ./pkg/client
	@cp -r ./tmp/github.com/keycloak/keycloak-operator/pkg/client/* ./pkg/client/
	@rm -rf ./tmp/github.com ./tmp/code-generator

.PHONY: test/goveralls
test/goveralls: test/coverage/prepare
	@echo "Preparing goveralls file"
	go get -u github.com/mattn/goveralls
	@echo "Running goveralls"
	@goveralls -v -coverprofile=cover-all.coverprofile -service=github
