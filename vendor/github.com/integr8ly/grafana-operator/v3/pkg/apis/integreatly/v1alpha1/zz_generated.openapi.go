// +build !ignore_autogenerated

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1alpha1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"./pkg/apis/integreatly/v1alpha1.Grafana":                 schema_pkg_apis_integreatly_v1alpha1_Grafana(ref),
		"./pkg/apis/integreatly/v1alpha1.GrafanaDashboard":        schema_pkg_apis_integreatly_v1alpha1_GrafanaDashboard(ref),
		"./pkg/apis/integreatly/v1alpha1.GrafanaDataSource":       schema_pkg_apis_integreatly_v1alpha1_GrafanaDataSource(ref),
		"./pkg/apis/integreatly/v1alpha1.GrafanaDataSourceSpec":   schema_pkg_apis_integreatly_v1alpha1_GrafanaDataSourceSpec(ref),
		"./pkg/apis/integreatly/v1alpha1.GrafanaDataSourceStatus": schema_pkg_apis_integreatly_v1alpha1_GrafanaDataSourceStatus(ref),
		"./pkg/apis/integreatly/v1alpha1.GrafanaSpec":             schema_pkg_apis_integreatly_v1alpha1_GrafanaSpec(ref),
		"./pkg/apis/integreatly/v1alpha1.GrafanaStatus":           schema_pkg_apis_integreatly_v1alpha1_GrafanaStatus(ref),
	}
}

func schema_pkg_apis_integreatly_v1alpha1_Grafana(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Grafana is the Schema for the grafanas API",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/integreatly/v1alpha1.GrafanaSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/integreatly/v1alpha1.GrafanaStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"./pkg/apis/integreatly/v1alpha1.GrafanaSpec", "./pkg/apis/integreatly/v1alpha1.GrafanaStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_integreatly_v1alpha1_GrafanaDashboard(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "GrafanaDashboard is the Schema for the grafanadashboards API",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/integreatly/v1alpha1.GrafanaDashboardSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/integreatly/v1alpha1.GrafanaDashboardStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"./pkg/apis/integreatly/v1alpha1.GrafanaDashboardSpec", "./pkg/apis/integreatly/v1alpha1.GrafanaDashboardStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_integreatly_v1alpha1_GrafanaDataSource(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "GrafanaDataSource is the Schema for the grafanadatasources API",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/integreatly/v1alpha1.GrafanaDataSourceSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/integreatly/v1alpha1.GrafanaDataSourceStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"./pkg/apis/integreatly/v1alpha1.GrafanaDataSourceSpec", "./pkg/apis/integreatly/v1alpha1.GrafanaDataSourceStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_integreatly_v1alpha1_GrafanaDataSourceSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "GrafanaDataSourceSpec defines the desired state of GrafanaDataSource",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"datasources": {
						SchemaProps: spec.SchemaProps{
							Description: "INSERT ADDITIONAL SPEC FIELDS - desired state of cluster Important: Run \"operator-sdk generate k8s\" to regenerate code after modifying this file Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("./pkg/apis/integreatly/v1alpha1.GrafanaDataSourceFields"),
									},
								},
							},
						},
					},
					"name": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
				},
				Required: []string{"datasources", "name"},
			},
		},
		Dependencies: []string{
			"./pkg/apis/integreatly/v1alpha1.GrafanaDataSourceFields"},
	}
}

func schema_pkg_apis_integreatly_v1alpha1_GrafanaDataSourceStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "GrafanaDataSourceStatus defines the observed state of GrafanaDataSource",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"phase": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"message": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
				},
				Required: []string{"phase", "message"},
			},
		},
	}
}

func schema_pkg_apis_integreatly_v1alpha1_GrafanaSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "GrafanaSpec defines the desired state of Grafana",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"config": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/integreatly/v1alpha1.GrafanaConfig"),
						},
					},
					"containers": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("k8s.io/api/core/v1.Container"),
									},
								},
							},
						},
					},
					"dashboardLabelSelector": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector"),
									},
								},
							},
						},
					},
					"ingress": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/integreatly/v1alpha1.GrafanaIngress"),
						},
					},
					"secrets": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
						},
					},
					"configMaps": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
						},
					},
					"service": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/integreatly/v1alpha1.GrafanaService"),
						},
					},
					"deployment": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/integreatly/v1alpha1.GrafanaDeployment"),
						},
					},
					"resources": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/api/core/v1.ResourceRequirements"),
						},
					},
					"serviceAccount": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/integreatly/v1alpha1.GrafanaServiceAccount"),
						},
					},
					"client": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/integreatly/v1alpha1.GrafanaClient"),
						},
					},
					"compat": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/integreatly/v1alpha1.GrafanaCompat"),
						},
					},
					"dashboardNamespaceSelector": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector"),
						},
					},
				},
				Required: []string{"config", "compat"},
			},
		},
		Dependencies: []string{
			"./pkg/apis/integreatly/v1alpha1.GrafanaClient", "./pkg/apis/integreatly/v1alpha1.GrafanaCompat", "./pkg/apis/integreatly/v1alpha1.GrafanaConfig", "./pkg/apis/integreatly/v1alpha1.GrafanaDeployment", "./pkg/apis/integreatly/v1alpha1.GrafanaIngress", "./pkg/apis/integreatly/v1alpha1.GrafanaService", "./pkg/apis/integreatly/v1alpha1.GrafanaServiceAccount", "k8s.io/api/core/v1.Container", "k8s.io/api/core/v1.ResourceRequirements", "k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector"},
	}
}

func schema_pkg_apis_integreatly_v1alpha1_GrafanaStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "GrafanaStatus defines the observed state of Grafana",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"phase": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"message": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"dashboards": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type: []string{"array"},
										Items: &spec.SchemaOrArray{
											Schema: &spec.Schema{
												SchemaProps: spec.SchemaProps{
													Ref: ref("./pkg/apis/integreatly/v1alpha1.GrafanaDashboardRef"),
												},
											},
										},
									},
								},
							},
						},
					},
					"installedPlugins": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("./pkg/apis/integreatly/v1alpha1.GrafanaPlugin"),
									},
								},
							},
						},
					},
					"failedPlugins": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("./pkg/apis/integreatly/v1alpha1.GrafanaPlugin"),
									},
								},
							},
						},
					},
				},
				Required: []string{"phase", "message", "dashboards", "installedPlugins", "failedPlugins"},
			},
		},
		Dependencies: []string{
			"./pkg/apis/integreatly/v1alpha1.GrafanaDashboardRef", "./pkg/apis/integreatly/v1alpha1.GrafanaPlugin"},
	}
}
