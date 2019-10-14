module github.com/keycloak/keycloak-operator

require (
	cloud.google.com/go v0.46.3 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/coreos/prometheus-operator v0.29.0
	github.com/go-openapi/spec v0.19.0
	github.com/gogo/protobuf v1.3.0 // indirect
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.11.2 // indirect
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/integr8ly/grafana-operator v1.4.0
	github.com/mattn/goveralls v0.0.3 // indirect
	github.com/onsi/ginkgo v1.10.1 // indirect
	github.com/onsi/gomega v1.7.0 // indirect
	github.com/openshift/api v3.9.0+incompatible
	github.com/operator-framework/operator-sdk v0.10.1-0.20190911145116-334c667503d0
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.1.0 // indirect
	github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4 // indirect
	github.com/prometheus/common v0.7.0 // indirect
	github.com/prometheus/procfs v0.0.5 // indirect
	github.com/rogpeppe/go-internal v1.4.0 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.4.0
	go.opencensus.io v0.22.1 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/zap v1.10.0 // indirect
	golang.org/x/crypto v0.0.0-20190926114937-fa1a29108794 // indirect
	golang.org/x/net v0.0.0-20190926025831-c00fd9afed17 // indirect
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e // indirect
	golang.org/x/sys v0.0.0-20190924154521-2837fb4f24fe // indirect
	golang.org/x/time v0.0.0-20190921001708-c4c64cad1fd0 // indirect
	golang.org/x/tools v0.0.0-20191012152004-8de300cfc20a // indirect
	google.golang.org/api v0.10.0 // indirect
	google.golang.org/appengine v1.6.3 // indirect
	google.golang.org/genproto v0.0.0-20190925194540-b8fbc687dcfb // indirect
	google.golang.org/grpc v1.24.0 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	k8s.io/api v0.0.0-20190612125737-db0771252981
	k8s.io/apimachinery v0.0.0-20190612125636-6a5db36e93ad
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/kube-openapi v0.0.0-20190603182131-db7b694dc208
	sigs.k8s.io/controller-runtime v0.1.12
	sigs.k8s.io/controller-tools v0.1.10
)

// Pinned to kubernetes-1.13.4
replace (
	k8s.io/api => k8s.io/api v0.0.0-20190222213804-5cb15d344471
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190228180357-d002e88f6236
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190221213512-86fb29eff628
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190228174230-b40b2a5939e4
)

replace (
	github.com/coreos/prometheus-operator => github.com/coreos/prometheus-operator v0.29.0
	// Pinned to v2.9.2 (kubernetes-1.13.1) so https://proxy.golang.org can
	// resolve it correctly.
	github.com/prometheus/prometheus => github.com/prometheus/prometheus v0.0.0-20190424153033-d3245f150225
	k8s.io/kube-state-metrics => k8s.io/kube-state-metrics v1.6.0
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.1.12
	sigs.k8s.io/controller-tools => sigs.k8s.io/controller-tools v0.1.11-0.20190411181648-9d55346c2bde
)

replace github.com/operator-framework/operator-sdk => github.com/operator-framework/operator-sdk v0.10.0

go 1.13
