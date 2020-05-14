package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	keycloakv1alpha1 "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/stretchr/testify/assert"

	"k8s.io/apimachinery/pkg/types"

	framework "github.com/operator-framework/operator-sdk/pkg/test"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	dynclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Condition func(t *testing.T, c kubernetes.Interface) error

func WaitForCondition(t *testing.T, c kubernetes.Interface, cond Condition) error {
	t.Logf("waiting up to %v for condition", pollTimeout)
	var err error
	for start := time.Now(); time.Since(start) < pollTimeout; time.Sleep(pollRetryInterval) {
		err = cond(t, c)
		if err == nil {
			return nil
		}
	}
	return err
}

// Stolen from https://github.com/kubernetes/kubernetes/blob/master/test/e2e/framework/util.go
// Then rewritten to use internal condition statements.
func WaitForStatefulSetReplicasReady(t *testing.T, c kubernetes.Interface, statefulSetName, ns string) error {
	t.Logf("waiting up to %v for StatefulSet %s to have all replicas ready", pollTimeout, statefulSetName)
	return WaitForCondition(t, c, func(t *testing.T, c kubernetes.Interface) error {
		sts, err := c.AppsV1().StatefulSets(ns).Get(statefulSetName, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("get StatefulSet %s failed, ignoring for %v: %v", statefulSetName, pollRetryInterval, err)
		}
		if sts.Status.ReadyReplicas == *sts.Spec.Replicas {
			t.Logf("all %d replicas of StatefulSet %s are ready.", sts.Status.ReadyReplicas, statefulSetName)
			return nil
		}
		return fmt.Errorf("statefulSet %s found but there are %d ready replicas and %d total replicas", statefulSetName, sts.Status.ReadyReplicas, *sts.Spec.Replicas)
	})
}

func WaitForPersistentVolumeClaimCreated(t *testing.T, c kubernetes.Interface, persistentVolumeClaimName, ns string) error {
	t.Logf("waiting up to %v for PersistentVolumeClaim %s to be created", pollTimeout, persistentVolumeClaimName)
	return WaitForCondition(t, c, func(t *testing.T, c kubernetes.Interface) error {
		pvc, err := c.CoreV1().PersistentVolumeClaims(ns).Get(persistentVolumeClaimName, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("get PersistentVolumeClaim %s failed, ignoring for %v: %v", persistentVolumeClaimName, pollRetryInterval, err)
		}
		if pvc.Status.Phase == "Bound" {
			t.Logf("PersistentVolumeClaim is bound")
			return nil
		}
		return fmt.Errorf("persistentVolumeClaim %s found but is not bound", persistentVolumeClaimName)
	})
}

func WaitForRealmToBeReady(t *testing.T, framework *framework.Framework, namespace string) error {
	keycloakRealmCR := &keycloakv1alpha1.KeycloakRealm{}

	return WaitForCondition(t, framework.KubeClient, func(t *testing.T, c kubernetes.Interface) error {
		err := GetNamespacedObject(framework, namespace, testKeycloakRealmCRName, keycloakRealmCR)
		if err != nil {
			return err
		}

		if !keycloakRealmCR.Status.Ready {
			keycloakRealmCRParsed, err := json.Marshal(keycloakRealmCR)
			if err != nil {
				return err
			}

			return fmt.Errorf("KeycloakRealm is not ready \nCurrent CR value: %s", string(keycloakRealmCRParsed))
		}

		return nil
	})
}

func WaitForClientToBeReady(t *testing.T, framework *framework.Framework, namespace string) error {
	keycloakClientCR := &keycloakv1alpha1.KeycloakClient{}

	return WaitForCondition(t, framework.KubeClient, func(t *testing.T, c kubernetes.Interface) error {
		err := GetNamespacedObject(framework, namespace, testKeycloakClientCRName, keycloakClientCR)
		if err != nil {
			return err
		}

		if !keycloakClientCR.Status.Ready {
			keycloakRealmCRParsed, err := json.Marshal(keycloakClientCR)
			if err != nil {
				return err
			}

			return fmt.Errorf("KeycloakClient is not ready \nCurrent CR value: %s", string(keycloakRealmCRParsed))
		}

		return nil
	})
}

func WaitForUserToBeReady(t *testing.T, framework *framework.Framework, namespace string) error {
	keycloakUserCR := &keycloakv1alpha1.KeycloakUser{}

	return WaitForCondition(t, framework.KubeClient, func(t *testing.T, c kubernetes.Interface) error {
		err := GetNamespacedObject(framework, namespace, testKeycloakUserCRName, keycloakUserCR)
		if err != nil {
			return err
		}

		if keycloakUserCR.Status.Phase != keycloakv1alpha1.UserPhaseReconciled {
			keycloakRealmCRParsed, err := json.Marshal(keycloakUserCR)
			if err != nil {
				return err
			}

			return fmt.Errorf("KeycloakRealm is not ready \nCurrent CR value: %s", string(keycloakRealmCRParsed))
		}

		return nil
	})
}

func WaitForSuccessResponseToContain(t *testing.T, framework *framework.Framework, url string, expectedString string) error {
	return WaitForCondition(t, framework.KubeClient, func(t *testing.T, c kubernetes.Interface) error {
		response, err := http.Get(url)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		if response.StatusCode != 200 {
			return fmt.Errorf("invalid response from url %s (%v)", url, response.Status)
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		responseString := string(responseData)

		assert.Contains(t, responseString, expectedString)

		return nil
	})
}

func Create(f *framework.Framework, obj runtime.Object, ctx *framework.TestCtx) error {
	return f.Client.Create(context.TODO(), obj, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
}

func Get(f *framework.Framework, key dynclient.ObjectKey, obj runtime.Object) error {
	return f.Client.Get(context.TODO(), key, obj)
}

func GetNamespacedObject(f *framework.Framework, namespace string, objectName string, outputObject runtime.Object) error {
	key := types.NamespacedName{
		Namespace: namespace,
		Name:      objectName,
	}

	return Get(f, key, outputObject)
}

func Update(f *framework.Framework, obj runtime.Object) error {
	return f.Client.Update(context.TODO(), obj)
}

func Delete(f *framework.Framework, obj runtime.Object) error {
	return f.Client.Delete(context.TODO(), obj)
}

func CreateLabel(namespace string) map[string]string {
	return map[string]string{"app": "keycloak-in-" + namespace}
}
