package e2e

import (
	"context"
	"fmt"
	"testing"
	"time"

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

func Create(f *framework.Framework, obj runtime.Object, ctx *framework.TestCtx) error {
	return f.Client.Create(context.TODO(), obj, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
}

func Get(f *framework.Framework, key dynclient.ObjectKey, obj runtime.Object) error {
	return f.Client.Get(context.TODO(), key, obj)
}

func Delete(f *framework.Framework, obj runtime.Object) error {
	return f.Client.Delete(context.TODO(), obj)
}

func CreateLabel(namespace string) map[string]string {
	return map[string]string{"app": "keycloak-in-" + namespace}
}
