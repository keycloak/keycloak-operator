package e2e

import (
	"github.com/keycloak/keycloak-operator/pkg/k8sutil"
	routev1 "github.com/openshift/api/route/v1"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// workaround for https://github.com/kubernetes/minikube/issues/3129
func doWorkaroundIfNecessary(f *framework.Framework, ctx *framework.Context, namespace string) error {
	resourceExists, _ := k8sutil.ResourceExists(f.KubeClient.Discovery(), routev1.SchemeGroupVersion.String(), "Route")

	if !resourceExists {
		// We are not on Openshift, so we need to workaround the issue

		postqresqlPVC := &v1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "keycloak-postgresql-claim",
				Labels:    map[string]string{"app": "keycloak"},
				Namespace: namespace,
			},
			Spec: v1.PersistentVolumeClaimSpec{
				AccessModes: []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
				Resources: v1.ResourceRequirements{
					Requests: v1.ResourceList{
						v1.ResourceStorage: resource.MustParse("1Gi"),
					}},
			},
		}

		err := Create(f, postqresqlPVC, ctx)
		if err != nil {
			return err
		}

		backupPVC := &v1.PersistentVolumeClaim{
			ObjectMeta: v12.ObjectMeta{
				Name: "keycloak-backup-keycloak-test",
				Labels: map[string]string{
					"app":       "keycloak",
					"component": "database-backup",
				},
				Namespace: namespace,
			},
			Spec: v1.PersistentVolumeClaimSpec{
				AccessModes: []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
				Resources: v1.ResourceRequirements{
					Requests: v1.ResourceList{
						v1.ResourceStorage: resource.MustParse("1Gi"),
					}},
			},
		}

		err = Create(f, backupPVC, ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
