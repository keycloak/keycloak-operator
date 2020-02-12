package model

import (
	kc "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"k8s.io/api/extensions/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetHost(cr *kc.Keycloak) string {
	if cr.Spec.Ingress == nil {
		return ""
	}
	return cr.Spec.externalAccess.Hostname
}

func GetPath(cr *cr.Keycloak) string {
	if cr.Spec.externalAccess == nil {
		return "/"
	}
	return cr.Spec.externalAccess.Path
}

func GetIngressLabels(cr *cr.Keycloak) map[string]string {
	if cr.Spec.externalAccess == nil {
		return nil
	}
	return cr.Spec.externalAccess.Labels
}

func GetIngressAnnotations(cr *kc.Keycloak, existing map[string]string) map[string]string {
	if cr.Spec.externalAccess == nil {
		return existing
	}
	return MergeAnnotations(cr.Spec.externalAccess.Annotations, existing)
}

func getIngressTLS(cr *kc.Keycloak) []v1beta1.IngressTLS {
	if cr.Spec.externalAccess == nil {
		return nil
	}

	if cr.Spec.externalAccess.TLSEnabled {
		return []v1beta1.IngressTLS{
			{
				Hosts:      []string{cr.Spec.externalAccess.Hostname},
				SecretName: cr.Spec.externalAccess.TLSSecretName,
			},
		}
	}
	return nil
}

func getIngressSpec(cr *v1alpha1.Grafana) v1beta1.IngressSpec {
  return v1beta1.IngressSpec{
  TLS: getIngressTLS(cr),
  Rules: []v1beta1.IngressRule{
    {
      Host: GetHost(cr),
      IngressRuleValue: v1beta1.IngressRuleValue{
        HTTP: &v1beta1.HTTPIngressRuleValue{
          Paths: []v1beta1.HTTPIngressPath{
            {
              Path: GetPath(cr),
              Backend: v1beta1.IngressBackend{
                ServiceName: ApplicationName,
                ServicePort: intstr.FromInt(KeycloakServicePort),
              },
            },
          },
        },
      },
    },
  },
}
}

func KeycloakIngress(cr *kc.Keycloak) *v1beta1.Ingress {

  return &v1beta1.Ingress{
		ObjectMeta: v1.ObjectMeta{
			Name:        ApplicationName,
			Namespace:   cr.Namespace,
			Labels:      GetIngressLabels(cr),
			Annotations: GetIngressAnnotations(cr, nil),
		},
		Spec: getIngressSpec(cr),
	}
}


func KeycloakIngressReconciled(cr *kc.Keycloak, currentState *v1beta1.Ingress) *v1beta1.Ingress {
  reconciled := currentState.DeepCopy()
	reconciled.Labels = GetIngressLabels(cr)
	reconciled.Annotations = GetIngressAnnotations(cr, currentState.Annotations)
	reconciled.Spec = getIngressSpec(cr)
	return reconciled
}

func KeycloakIngressSelector(cr *kc.Keycloak) client.ObjectKey {
	return client.ObjectKey{
		Name:      ApplicationName,
		Namespace: cr.Namespace,
	}
}
