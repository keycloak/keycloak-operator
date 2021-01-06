package model

import (
	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func PrometheusRule(cr *v1alpha1.Keycloak) *monitoringv1.PrometheusRule {
	rules := []monitoringv1.Rule{{
		Alert: "KeycloakJavaNonHeapThresholdExceeded",
		Annotations: map[string]string{
			"message": `{{ printf "%0.0f" $value }}% nonheap usage of {{ $labels.area }} in pod {{ $labels.pod }}, namespace {{ $labels.namespace }}.`,
		},
		Expr: intstr.FromString(`100 * jvm_memory_bytes_used{area="nonheap",namespace="` + cr.Namespace + `"} / jvm_memory_bytes_max{area="nonheap",namespace="` + cr.Namespace + `"} > 90`),
		For:  "1m",
		Labels: map[string]string{
			"severity": "warning",
		},
	}, {
		Alert: "KeycloakJavaGCTimePerMinuteScavenge",
		Annotations: map[string]string{
			"message": `Amount of time per minute spent on garbage collection of {{ $labels.area }} in pod {{ $labels.pod }}, namespace {{ $labels.namespace }} exceeds 90%. This could indicate that the available heap memory is insufficient.`,
		},
		Expr: intstr.FromString(`increase(jvm_gc_collection_seconds_sum{gc="PS Scavenge",namespace="` + cr.Namespace + `"}[1m]) > 1 * 60 * 0.9`),
		For:  "1m",
		Labels: map[string]string{
			"severity": "warning",
		},
	}, {
		Alert: "KeycloakJavaGCTimePerMinuteMarkSweep",
		Annotations: map[string]string{
			"message": `Amount of time per minute spent on garbage collection of {{ $labels.area }} in pod {{ $labels.pod }}, namespace {{ $labels.namespace }} exceeds 90%. This could indicate that the available heap memory is insufficient.`,
		},
		Expr: intstr.FromString(`increase(jvm_gc_collection_seconds_sum{gc="PS MarkSweep",namespace="` + cr.Namespace + `"}[1m]) > 1 * 60 * 0.9`),
		For:  "1m",
		Labels: map[string]string{
			"severity": "warning",
		},
	}, {
		Alert: "KeycloakJavaDeadlockedThreads",
		Annotations: map[string]string{
			"message": `Number of threads in deadlock state of {{ $labels.area }} in pod {{ $labels.pod }}, namespace {{ $labels.namespace }}`,
		},
		Expr: intstr.FromString(`jvm_threads_deadlocked{namespace="` + cr.Namespace + `"} > 0`),
		For:  "1m",
		Labels: map[string]string{
			"severity": "warning",
		},
	}, {
		Alert: "KeycloakLoginFailedThresholdExceeded",
		Annotations: map[string]string{
			"message": `More than 50 failed login attempts for realm {{ $labels.realm }}, provider {{ $labels.provider }}, namespace {{ $labels.namespace }} over the last 5 minutes. (Rate of {{ printf "%0f" $value }})`,
		},
		Expr: intstr.FromString(`rate(keycloak_failed_login_attempts{namespace="` + cr.Namespace + `"}[5m]) * 300 > 50`),
		For:  "5m",
		Labels: map[string]string{
			"severity": "warning",
		},
	}, {
		Alert: "KeycloakInstanceNotAvailable",
		Annotations: map[string]string{
			"message": `Keycloak instance in namespace {{ $labels.namespace }} has not been available for the last 5 minutes.`,
		},
		Expr: intstr.FromString(`(1 - absent(kube_pod_status_ready{namespace="` + cr.Namespace + `", condition="true"} * on (pod) group_left (label_component) kube_pod_labels{label_component="` + KeycloakDeploymentComponent + `", namespace="` + cr.Namespace + `"})) == 0`),
		For:  "5m",
		Labels: map[string]string{
			"severity": "critical",
		},
	}, {
		Alert: "KeycloakAPIRequestDuration90PercThresholdExceeded",
		Annotations: map[string]string{
			"message": `More than 10% the RH SSO API endpoints in namespace {{ $labels.namespace }} are taking longer than 1s for the last 5 minutes.`,
		},
		Expr: intstr.FromString(`(sum(rate(keycloak_request_duration_bucket{le="1000.0", namespace="` + cr.Namespace + `"}[5m])) by (job) / sum(rate(keycloak_request_duration_count{namespace="` + cr.Namespace + `"}[5m])) by (job)) < 0.90`),
		For:  "5m",
		Labels: map[string]string{
			"severity": "warning",
		},
	}, {
		Alert: "KeycloakAPIRequestDuration99.5PercThresholdExceeded",
		Annotations: map[string]string{
			"message": `More than 0.5% of the RH SSO API endpoints in namespace {{ $labels.namespace }} are taking longer than 10s for the last 5 minutes.`,
		},
		Expr: intstr.FromString(`(sum(rate(keycloak_request_duration_bucket{le="10000.0", namespace="` + cr.Namespace + `"}[5m])) by (job) / sum(rate(keycloak_request_duration_count{namespace="` + cr.Namespace + `"}[5m])) by (job)) < 0.995`),
		For:  "5m",
		Labels: map[string]string{
			"severity": "warning",
		},
	}}

	if !cr.Spec.ExternalDatabase.Enabled {
		rules = append(rules, monitoringv1.Rule{
			Alert: "KeycloakDatabaseNotAvailable",
			Annotations: map[string]string{
				"message": `RH SSO database in namespace {{ $labels.namespace }} is not available for the last 5 minutes.`,
			},
			Expr: intstr.FromString(`(1 - absent(kube_pod_status_ready{namespace="` + cr.Namespace + `", condition="true"} * on (pod) group_left (label_component) kube_pod_labels{label_component="` + PostgresqlDeploymentComponent + `", namespace="` + cr.Namespace + `"})) == 0 `),
			For:  "5m",
			Labels: map[string]string{
				"severity": "critical",
			},
		})
	}

	return &monitoringv1.PrometheusRule{
		ObjectMeta: v12.ObjectMeta{
			Name:      ApplicationName,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"monitoring-key": MonitoringKey,
				"prometheus":     "application-monitoring",
				"role":           "alert-rules",
			},
		},
		Spec: monitoringv1.PrometheusRuleSpec{
			Groups: []monitoringv1.RuleGroup{{
				Name:  "general.rules",
				Rules: rules,
			}},
		},
	}
}

func PrometheusRuleSelector(cr *v1alpha1.Keycloak) client.ObjectKey {
	return client.ObjectKey{
		Name:      ApplicationName,
		Namespace: cr.Namespace,
	}
}
