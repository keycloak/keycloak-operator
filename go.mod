module github.com/keycloak/keycloak-operator

require (
	github.com/coreos/prometheus-operator v0.40.0
	github.com/go-openapi/spec v0.19.7
	github.com/integr8ly/grafana-operator/v3 v3.6.0
	github.com/json-iterator/go v1.1.9
	github.com/openshift/api v3.9.0+incompatible
	github.com/operator-framework/operator-sdk v0.18.2
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.5.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.6.1
	k8s.io/api v0.18.3
	k8s.io/apiextensions-apiserver v0.18.2
	k8s.io/apimachinery v0.18.3
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/kube-openapi v0.0.0-20200410145947-61e04a5be9a6
	k8s.io/utils v0.0.0-20200414100711-2df71ebbae66
	sigs.k8s.io/controller-runtime v0.6.0

)

// Pinned to kubernetes-1.18.2
replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.2.0+incompatible
	github.com/operator-framework/operator-sdk => github.com/operator-framework/operator-sdk v0.18.2
	k8s.io/api => k8s.io/api v0.18.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.2
)

replace k8s.io/client-go => k8s.io/client-go v0.18.2

go 1.13
