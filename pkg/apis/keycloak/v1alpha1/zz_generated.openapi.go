//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by openapi-gen. DO NOT EDIT.

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1alpha1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"./pkg/apis/keycloak/v1alpha1.Keycloak":             schema_pkg_apis_keycloak_v1alpha1_Keycloak(ref),
		"./pkg/apis/keycloak/v1alpha1.KeycloakAWSSpec":      schema_pkg_apis_keycloak_v1alpha1_KeycloakAWSSpec(ref),
		"./pkg/apis/keycloak/v1alpha1.KeycloakBackup":       schema_pkg_apis_keycloak_v1alpha1_KeycloakBackup(ref),
		"./pkg/apis/keycloak/v1alpha1.KeycloakBackupSpec":   schema_pkg_apis_keycloak_v1alpha1_KeycloakBackupSpec(ref),
		"./pkg/apis/keycloak/v1alpha1.KeycloakBackupStatus": schema_pkg_apis_keycloak_v1alpha1_KeycloakBackupStatus(ref),
		"./pkg/apis/keycloak/v1alpha1.KeycloakClient":       schema_pkg_apis_keycloak_v1alpha1_KeycloakClient(ref),
		"./pkg/apis/keycloak/v1alpha1.KeycloakClientSpec":   schema_pkg_apis_keycloak_v1alpha1_KeycloakClientSpec(ref),
		"./pkg/apis/keycloak/v1alpha1.KeycloakClientStatus": schema_pkg_apis_keycloak_v1alpha1_KeycloakClientStatus(ref),
		"./pkg/apis/keycloak/v1alpha1.KeycloakRealm":        schema_pkg_apis_keycloak_v1alpha1_KeycloakRealm(ref),
		"./pkg/apis/keycloak/v1alpha1.KeycloakRealmSpec":    schema_pkg_apis_keycloak_v1alpha1_KeycloakRealmSpec(ref),
		"./pkg/apis/keycloak/v1alpha1.KeycloakRealmStatus":  schema_pkg_apis_keycloak_v1alpha1_KeycloakRealmStatus(ref),
		"./pkg/apis/keycloak/v1alpha1.KeycloakSpec":         schema_pkg_apis_keycloak_v1alpha1_KeycloakSpec(ref),
		"./pkg/apis/keycloak/v1alpha1.KeycloakStatus":       schema_pkg_apis_keycloak_v1alpha1_KeycloakStatus(ref),
		"./pkg/apis/keycloak/v1alpha1.KeycloakUser":         schema_pkg_apis_keycloak_v1alpha1_KeycloakUser(ref),
		"./pkg/apis/keycloak/v1alpha1.KeycloakUserSpec":     schema_pkg_apis_keycloak_v1alpha1_KeycloakUserSpec(ref),
		"./pkg/apis/keycloak/v1alpha1.KeycloakUserStatus":   schema_pkg_apis_keycloak_v1alpha1_KeycloakUserStatus(ref),
	}
}

func schema_pkg_apis_keycloak_v1alpha1_Keycloak(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Keycloak is the Schema for the keycloaks API.",
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
							Default: map[string]interface{}{},
							Ref:     ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Default: map[string]interface{}{},
							Ref:     ref("./pkg/apis/keycloak/v1alpha1.KeycloakSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Default: map[string]interface{}{},
							Ref:     ref("./pkg/apis/keycloak/v1alpha1.KeycloakStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"./pkg/apis/keycloak/v1alpha1.KeycloakSpec", "./pkg/apis/keycloak/v1alpha1.KeycloakStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_keycloak_v1alpha1_KeycloakAWSSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "KeycloakAWSSpec defines the desired state of KeycloakBackupSpec.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"encryptionKeySecretName": {
						SchemaProps: spec.SchemaProps{
							Description: "If provided, the database backup will be encrypted. Provides a secret name used for encrypting database data. The secret needs to be in the following form:\n\n    apiVersion: v1\n    kind: Secret\n    metadata:\n      name: <Secret name>\n    type: Opaque\n    stringData:\n      GPG_PUBLIC_KEY: <GPG Public Key>\n      GPG_TRUST_MODEL: <GPG Trust Model>\n      GPG_RECIPIENT: <GPG Recipient>\n\nFor more information, please refer to the Operator documentation.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"credentialsSecretName": {
						SchemaProps: spec.SchemaProps{
							Description: "Provides a secret name used for connecting to AWS S3 Service. The secret needs to be in the following form:\n\n    apiVersion: v1\n    kind: Secret\n    metadata:\n      name: <Secret name>\n    type: Opaque\n    stringData:\n      AWS_S3_BUCKET_NAME: <S3 Bucket Name>\n      AWS_ACCESS_KEY_ID: <AWS Access Key ID>\n      AWS_SECRET_ACCESS_KEY: <AWS Secret Key>\n\nFor more information, please refer to the Operator documentation.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"schedule": {
						SchemaProps: spec.SchemaProps{
							Description: "If specified, it will be used as a schedule for creating a CronJob.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
			},
		},
	}
}

func schema_pkg_apis_keycloak_v1alpha1_KeycloakBackup(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "KeycloakBackup is the Schema for the keycloakbackups API.",
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
							Default: map[string]interface{}{},
							Ref:     ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Default: map[string]interface{}{},
							Ref:     ref("./pkg/apis/keycloak/v1alpha1.KeycloakBackupSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Default: map[string]interface{}{},
							Ref:     ref("./pkg/apis/keycloak/v1alpha1.KeycloakBackupStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"./pkg/apis/keycloak/v1alpha1.KeycloakBackupSpec", "./pkg/apis/keycloak/v1alpha1.KeycloakBackupStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_keycloak_v1alpha1_KeycloakBackupSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "KeycloakBackupSpec defines the desired state of KeycloakBackup.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"restore": {
						SchemaProps: spec.SchemaProps{
							Description: "Controls automatic restore behavior. Currently not implemented.\n\nIn the future this will be used to trigger automatic restore for a given KeycloakBackup. Each backup will correspond to a single snapshot of the database (stored either in a Persistent Volume or AWS). If a user wants to restore it, all he/she needs to do is to change this flag to true. Potentially, it will be possible to restore a single backup multiple times.",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"aws": {
						SchemaProps: spec.SchemaProps{
							Description: "If provided, an automatic database backup will be created on AWS S3 instead of a local Persistent Volume. If this property is not provided - a local Persistent Volume backup will be chosen.",
							Default:     map[string]interface{}{},
							Ref:         ref("./pkg/apis/keycloak/v1alpha1.KeycloakAWSSpec"),
						},
					},
					"instanceSelector": {
						SchemaProps: spec.SchemaProps{
							Description: "Selector for looking up Keycloak Custom Resources.",
							Ref:         ref("k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector"),
						},
					},
					"storageClassName": {
						SchemaProps: spec.SchemaProps{
							Description: "Name of the StorageClass for Postgresql Backup Persistent Volume Claim",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
			},
		},
		Dependencies: []string{
			"./pkg/apis/keycloak/v1alpha1.KeycloakAWSSpec", "k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector"},
	}
}

func schema_pkg_apis_keycloak_v1alpha1_KeycloakBackupStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "KeycloakBackupStatus defines the observed state of KeycloakBackup.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"phase": {
						SchemaProps: spec.SchemaProps{
							Description: "Current phase of the operator.",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"message": {
						SchemaProps: spec.SchemaProps{
							Description: "Human-readable message indicating details about current operator phase or error.",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"ready": {
						SchemaProps: spec.SchemaProps{
							Description: "True if all resources are in a ready state and all work is done.",
							Default:     false,
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"secondaryResources": {
						SchemaProps: spec.SchemaProps{
							Description: "A map of all the secondary resources types and names created for this CR. e.g \"Deployment\": [ \"DeploymentName1\", \"DeploymentName2\" ]",
							Type:        []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type: []string{"array"},
										Items: &spec.SchemaOrArray{
											Schema: &spec.Schema{
												SchemaProps: spec.SchemaProps{
													Default: "",
													Type:    []string{"string"},
													Format:  "",
												},
											},
										},
									},
								},
							},
						},
					},
				},
				Required: []string{"phase", "message", "ready"},
			},
		},
	}
}

func schema_pkg_apis_keycloak_v1alpha1_KeycloakClient(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "KeycloakClient is the Schema for the keycloakclients API.",
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
							Default: map[string]interface{}{},
							Ref:     ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Default: map[string]interface{}{},
							Ref:     ref("./pkg/apis/keycloak/v1alpha1.KeycloakClientSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Default: map[string]interface{}{},
							Ref:     ref("./pkg/apis/keycloak/v1alpha1.KeycloakClientStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"./pkg/apis/keycloak/v1alpha1.KeycloakClientSpec", "./pkg/apis/keycloak/v1alpha1.KeycloakClientStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_keycloak_v1alpha1_KeycloakClientSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "KeycloakClientSpec defines the desired state of KeycloakClient.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"realmSelector": {
						SchemaProps: spec.SchemaProps{
							Description: "Selector for looking up KeycloakRealm Custom Resources.",
							Ref:         ref("k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector"),
						},
					},
					"client": {
						SchemaProps: spec.SchemaProps{
							Description: "Keycloak Client REST object.",
							Ref:         ref("./pkg/apis/keycloak/v1alpha1.KeycloakAPIClient"),
						},
					},
					"roles": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-list-map-keys": []interface{}{
									"name",
								},
								"x-kubernetes-list-type": "map",
							},
						},
						SchemaProps: spec.SchemaProps{
							Description: "Client Roles",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: map[string]interface{}{},
										Ref:     ref("./pkg/apis/keycloak/v1alpha1.RoleRepresentation"),
									},
								},
							},
						},
					},
					"scopeMappings": {
						SchemaProps: spec.SchemaProps{
							Description: "Scope Mappings",
							Ref:         ref("./pkg/apis/keycloak/v1alpha1.MappingsRepresentation"),
						},
					},
				},
				Required: []string{"realmSelector", "client"},
			},
		},
		Dependencies: []string{
			"./pkg/apis/keycloak/v1alpha1.KeycloakAPIClient", "./pkg/apis/keycloak/v1alpha1.MappingsRepresentation", "./pkg/apis/keycloak/v1alpha1.RoleRepresentation", "k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector"},
	}
}

func schema_pkg_apis_keycloak_v1alpha1_KeycloakClientStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "KeycloakClientStatus defines the observed state of KeycloakClient",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"phase": {
						SchemaProps: spec.SchemaProps{
							Description: "Current phase of the operator.",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"message": {
						SchemaProps: spec.SchemaProps{
							Description: "Human-readable message indicating details about current operator phase or error.",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"ready": {
						SchemaProps: spec.SchemaProps{
							Description: "True if all resources are in a ready state and all work is done.",
							Default:     false,
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"secondaryResources": {
						SchemaProps: spec.SchemaProps{
							Description: "A map of all the secondary resources types and names created for this CR. e.g \"Deployment\": [ \"DeploymentName1\", \"DeploymentName2\" ]",
							Type:        []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type: []string{"array"},
										Items: &spec.SchemaOrArray{
											Schema: &spec.Schema{
												SchemaProps: spec.SchemaProps{
													Default: "",
													Type:    []string{"string"},
													Format:  "",
												},
											},
										},
									},
								},
							},
						},
					},
				},
				Required: []string{"phase", "message", "ready"},
			},
		},
	}
}

func schema_pkg_apis_keycloak_v1alpha1_KeycloakRealm(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "KeycloakRealm is the Schema for the keycloakrealms API",
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
							Default: map[string]interface{}{},
							Ref:     ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Default: map[string]interface{}{},
							Ref:     ref("./pkg/apis/keycloak/v1alpha1.KeycloakRealmSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Default: map[string]interface{}{},
							Ref:     ref("./pkg/apis/keycloak/v1alpha1.KeycloakRealmStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"./pkg/apis/keycloak/v1alpha1.KeycloakRealmSpec", "./pkg/apis/keycloak/v1alpha1.KeycloakRealmStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_keycloak_v1alpha1_KeycloakRealmSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "KeycloakRealmSpec defines the desired state of KeycloakRealm.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"unmanaged": {
						SchemaProps: spec.SchemaProps{
							Description: "When set to true, this KeycloakRealm will be marked as unmanaged and not be managed by this operator. It can then be used for targeting purposes.",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"instanceSelector": {
						SchemaProps: spec.SchemaProps{
							Description: "Selector for looking up Keycloak Custom Resources.",
							Ref:         ref("k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector"),
						},
					},
					"realm": {
						SchemaProps: spec.SchemaProps{
							Description: "Keycloak Realm REST object.",
							Ref:         ref("./pkg/apis/keycloak/v1alpha1.KeycloakAPIRealm"),
						},
					},
					"realmOverrides": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-list-type": "atomic",
							},
						},
						SchemaProps: spec.SchemaProps{
							Description: "A list of overrides to the default Realm behavior.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("./pkg/apis/keycloak/v1alpha1.RedirectorIdentityProviderOverride"),
									},
								},
							},
						},
					},
				},
				Required: []string{"realm"},
			},
		},
		Dependencies: []string{
			"./pkg/apis/keycloak/v1alpha1.KeycloakAPIRealm", "./pkg/apis/keycloak/v1alpha1.RedirectorIdentityProviderOverride", "k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector"},
	}
}

func schema_pkg_apis_keycloak_v1alpha1_KeycloakRealmStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "KeycloakRealmStatus defines the observed state of KeycloakRealm",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"phase": {
						SchemaProps: spec.SchemaProps{
							Description: "Current phase of the operator.",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"message": {
						SchemaProps: spec.SchemaProps{
							Description: "Human-readable message indicating details about current operator phase or error.",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"ready": {
						SchemaProps: spec.SchemaProps{
							Description: "True if all resources are in a ready state and all work is done.",
							Default:     false,
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"secondaryResources": {
						SchemaProps: spec.SchemaProps{
							Description: "A map of all the secondary resources types and names created for this CR. e.g \"Deployment\": [ \"DeploymentName1\", \"DeploymentName2\" ]",
							Type:        []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type: []string{"array"},
										Items: &spec.SchemaOrArray{
											Schema: &spec.Schema{
												SchemaProps: spec.SchemaProps{
													Default: "",
													Type:    []string{"string"},
													Format:  "",
												},
											},
										},
									},
								},
							},
						},
					},
					"loginURL": {
						SchemaProps: spec.SchemaProps{
							Default: "",
							Type:    []string{"string"},
							Format:  "",
						},
					},
				},
				Required: []string{"phase", "message", "ready", "loginURL"},
			},
		},
	}
}

func schema_pkg_apis_keycloak_v1alpha1_KeycloakSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "KeycloakSpec defines the desired state of Keycloak.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"unmanaged": {
						SchemaProps: spec.SchemaProps{
							Description: "When set to true, this Keycloak will be marked as unmanaged and will not be managed by this operator. It can then be used for targeting purposes.",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"external": {
						SchemaProps: spec.SchemaProps{
							Description: "Contains configuration for external Keycloak instances. Unmanaged needs to be set to true to use this.",
							Default:     map[string]interface{}{},
							Ref:         ref("./pkg/apis/keycloak/v1alpha1.KeycloakExternal"),
						},
					},
					"extensions": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-list-type": "set",
							},
						},
						SchemaProps: spec.SchemaProps{
							Description: "A list of extensions, where each one is a URL to a JAR files that will be deployed in Keycloak.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: "",
										Type:    []string{"string"},
										Format:  "",
									},
								},
							},
						},
					},
					"instances": {
						SchemaProps: spec.SchemaProps{
							Description: "Number of Keycloak instances in HA mode. Default is 1.",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"externalAccess": {
						SchemaProps: spec.SchemaProps{
							Description: "Controls external Ingress/Route settings.",
							Default:     map[string]interface{}{},
							Ref:         ref("./pkg/apis/keycloak/v1alpha1.KeycloakExternalAccess"),
						},
					},
					"externalDatabase": {
						SchemaProps: spec.SchemaProps{
							Description: "Controls external database settings. Using an external database requires providing a secret containing credentials as well as connection details. Here's an example of such secret:\n\n    apiVersion: v1\n    kind: Secret\n    metadata:\n        name: keycloak-db-secret\n        namespace: keycloak\n    stringData:\n        POSTGRES_DATABASE: <Database Name>\n        POSTGRES_EXTERNAL_ADDRESS: <External Database IP or URL (resolvable by K8s)>\n        POSTGRES_EXTERNAL_PORT: <External Database Port>\n        # Strongly recommended to use <'Keycloak CR Name'-postgresql>\n        POSTGRES_HOST: <Database Service Name>\n        POSTGRES_PASSWORD: <Database Password>\n        # Required for AWS Backup functionality\n        POSTGRES_SUPERUSER: true\n        POSTGRES_USERNAME: <Database Username>\n     type: Opaque\n\nBoth POSTGRES_EXTERNAL_ADDRESS and POSTGRES_EXTERNAL_PORT are specifically required for creating connection to the external database. The secret name is created using the following convention:\n      <Custom Resource Name>-db-secret\n\nFor more information, please refer to the Operator documentation.",
							Default:     map[string]interface{}{},
							Ref:         ref("./pkg/apis/keycloak/v1alpha1.KeycloakExternalDatabase"),
						},
					},
					"profile": {
						SchemaProps: spec.SchemaProps{
							Description: "Profile used for controlling Operator behavior. Default is empty.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"podDisruptionBudget": {
						SchemaProps: spec.SchemaProps{
							Description: "Specify PodDisruptionBudget configuration.",
							Default:     map[string]interface{}{},
							Ref:         ref("./pkg/apis/keycloak/v1alpha1.PodDisruptionBudgetConfig"),
						},
					},
					"keycloakDeploymentSpec": {
						SchemaProps: spec.SchemaProps{
							Description: "Resources (Requests and Limits) and ImagePullPolicy for KeycloakDeployment.",
							Default:     map[string]interface{}{},
							Ref:         ref("./pkg/apis/keycloak/v1alpha1.KeycloakDeploymentSpec"),
						},
					},
					"postgresDeploymentSpec": {
						SchemaProps: spec.SchemaProps{
							Description: "Resources (Requests and Limits) and ImagePullPolicy for PostgresDeployment.",
							Default:     map[string]interface{}{},
							Ref:         ref("./pkg/apis/keycloak/v1alpha1.PostgresqlDeploymentSpec"),
						},
					},
					"migration": {
						SchemaProps: spec.SchemaProps{
							Description: "Specify Migration configuration",
							Default:     map[string]interface{}{},
							Ref:         ref("./pkg/apis/keycloak/v1alpha1.MigrateConfig"),
						},
					},
					"storageClassName": {
						SchemaProps: spec.SchemaProps{
							Description: "Name of the StorageClass for Postgresql Persistent Volume Claim",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"multiAvailablityZones": {
						SchemaProps: spec.SchemaProps{
							Description: "Specify PodAntiAffinity settings for Keycloak deployment in Multi AZ",
							Default:     map[string]interface{}{},
							Ref:         ref("./pkg/apis/keycloak/v1alpha1.MultiAvailablityZonesConfig"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"./pkg/apis/keycloak/v1alpha1.KeycloakDeploymentSpec", "./pkg/apis/keycloak/v1alpha1.KeycloakExternal", "./pkg/apis/keycloak/v1alpha1.KeycloakExternalAccess", "./pkg/apis/keycloak/v1alpha1.KeycloakExternalDatabase", "./pkg/apis/keycloak/v1alpha1.MigrateConfig", "./pkg/apis/keycloak/v1alpha1.MultiAvailablityZonesConfig", "./pkg/apis/keycloak/v1alpha1.PodDisruptionBudgetConfig", "./pkg/apis/keycloak/v1alpha1.PostgresqlDeploymentSpec"},
	}
}

func schema_pkg_apis_keycloak_v1alpha1_KeycloakStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "KeycloakStatus defines the observed state of Keycloak.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"phase": {
						SchemaProps: spec.SchemaProps{
							Description: "Current phase of the operator.",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"message": {
						SchemaProps: spec.SchemaProps{
							Description: "Human-readable message indicating details about current operator phase or error.",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"ready": {
						SchemaProps: spec.SchemaProps{
							Description: "True if all resources are in a ready state and all work is done.",
							Default:     false,
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"secondaryResources": {
						SchemaProps: spec.SchemaProps{
							Description: "A map of all the secondary resources types and names created for this CR. e.g \"Deployment\": [ \"DeploymentName1\", \"DeploymentName2\" ].",
							Type:        []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type: []string{"array"},
										Items: &spec.SchemaOrArray{
											Schema: &spec.Schema{
												SchemaProps: spec.SchemaProps{
													Default: "",
													Type:    []string{"string"},
													Format:  "",
												},
											},
										},
									},
								},
							},
						},
					},
					"version": {
						SchemaProps: spec.SchemaProps{
							Description: "Version of Keycloak or RHSSO running on the cluster.",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"internalURL": {
						SchemaProps: spec.SchemaProps{
							Description: "An internal URL (service name) to be used by the admin client.",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"externalURL": {
						SchemaProps: spec.SchemaProps{
							Description: "External URL for accessing Keycloak instance from outside the cluster. Is identical to external.URL if it's specified, otherwise is computed (e.g. from Ingress).",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"credentialSecret": {
						SchemaProps: spec.SchemaProps{
							Description: "The secret where the admin credentials are to be found.",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
				Required: []string{"phase", "message", "ready", "version", "internalURL", "credentialSecret"},
			},
		},
	}
}

func schema_pkg_apis_keycloak_v1alpha1_KeycloakUser(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "KeycloakUser is the Schema for the keycloakusers API.",
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
							Default: map[string]interface{}{},
							Ref:     ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Default: map[string]interface{}{},
							Ref:     ref("./pkg/apis/keycloak/v1alpha1.KeycloakUserSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Default: map[string]interface{}{},
							Ref:     ref("./pkg/apis/keycloak/v1alpha1.KeycloakUserStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"./pkg/apis/keycloak/v1alpha1.KeycloakUserSpec", "./pkg/apis/keycloak/v1alpha1.KeycloakUserStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_keycloak_v1alpha1_KeycloakUserSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "KeycloakUserSpec defines the desired state of KeycloakUser.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"realmSelector": {
						SchemaProps: spec.SchemaProps{
							Description: "Selector for looking up KeycloakRealm Custom Resources.",
							Ref:         ref("k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector"),
						},
					},
					"user": {
						SchemaProps: spec.SchemaProps{
							Description: "Keycloak User REST object.",
							Default:     map[string]interface{}{},
							Ref:         ref("./pkg/apis/keycloak/v1alpha1.KeycloakAPIUser"),
						},
					},
				},
				Required: []string{"user"},
			},
		},
		Dependencies: []string{
			"./pkg/apis/keycloak/v1alpha1.KeycloakAPIUser", "k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector"},
	}
}

func schema_pkg_apis_keycloak_v1alpha1_KeycloakUserStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "KeycloakUserStatus defines the observed state of KeycloakUser.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"phase": {
						SchemaProps: spec.SchemaProps{
							Description: "Current phase of the operator.",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"message": {
						SchemaProps: spec.SchemaProps{
							Description: "Human-readable message indicating details about current operator phase or error.",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
				Required: []string{"phase", "message"},
			},
		},
	}
}
